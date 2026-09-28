package bot

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/YutoMaeda1209/hygge/model"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

const handlerTimeout = 30 * time.Second

// userPlaceholder in a lobby's NewChannelName is replaced with the joining user's display name.
const userPlaceholder = "{user}"

// onGuildCreate deletes rooms left empty while the bot was offline.
// The state cache already holds the guild's voice states when this runs.
func onGuildCreate(s *discordgo.Session, g *discordgo.GuildCreate) {
	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeout)
	defer cancel()

	channelIds, err := model.ListVoiceLobbyRoomChannelIds(ctx, g.ID)
	if err != nil {
		slog.Error("Failed to list voice lobby rooms.", "guildId", g.ID, "err", err)
		return
	}
	for _, channelId := range channelIds {
		deleteRoomIfEmpty(ctx, s, g.ID, channelId)
	}
}

func onVoiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	beforeChannelId := ""
	if v.BeforeUpdate != nil {
		beforeChannelId = v.BeforeUpdate.ChannelID
	}
	if beforeChannelId == v.ChannelID {
		// Mute, deafen and similar changes.
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeout)
	defer cancel()

	if beforeChannelId != "" {
		deleteRoomIfEmpty(ctx, s, v.GuildID, beforeChannelId)
	}
	if v.ChannelID != "" {
		handleLobbyJoin(ctx, s, v)
	}
}

// onChannelDelete forgets rooms deleted by someone other than the bot.
func onChannelDelete(s *discordgo.Session, c *discordgo.ChannelDelete) {
	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeout)
	defer cancel()

	if err := model.DeleteVoiceLobbyRoom(ctx, c.ID); err != nil {
		slog.Error("Failed to delete a voice lobby room record.", "channelId", c.ID, "err", err)
	}
}

// handleLobbyJoin creates a room next to the lobby the user joined and moves the user into it.
// It does nothing if the joined channel is not a lobby.
func handleLobbyJoin(ctx context.Context, s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	var lobby model.VoiceLobbyChannel
	err := lobby.FindByChannelId(ctx, v.ChannelID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return
	} else if err != nil {
		slog.Error("Failed to find a voice lobby channel.", "channelId", v.ChannelID, "err", err)
		return
	}

	lobbyChannel, err := s.State.Channel(v.ChannelID)
	if err != nil {
		lobbyChannel, err = s.Channel(v.ChannelID, discordgo.WithContext(ctx))
		if err != nil {
			slog.Error("Failed to get the voice lobby channel.", "channelId", v.ChannelID, "err", err)
			return
		}
	}

	room, err := s.GuildChannelCreateComplex(v.GuildID, discordgo.GuildChannelCreateData{
		Name:     roomName(lobby.NewChannelName, v.Member),
		Type:     discordgo.ChannelTypeGuildVoice,
		ParentID: lobbyChannel.ParentID,
	}, discordgo.WithContext(ctx))
	if err != nil {
		slog.Error("Failed to create a voice lobby room.", "guildId", v.GuildID, "err", err)
		return
	}

	if err := model.CreateVoiceLobbyRoom(ctx, lobby.Id, room.ID); err != nil {
		slog.Error("Failed to record a voice lobby room.", "channelId", room.ID, "err", err)
		if _, err := s.ChannelDelete(room.ID, discordgo.WithContext(ctx)); err != nil {
			slog.Error("Failed to delete an unrecorded voice lobby room.", "channelId", room.ID, "err", err)
		}
		return
	}

	if err := s.GuildMemberMove(v.GuildID, v.UserID, &room.ID, discordgo.WithContext(ctx)); err != nil {
		// The user may have left the lobby before being moved.
		slog.Warn("Failed to move a user into a voice lobby room.", "userId", v.UserID, "channelId", room.ID, "err", err)
		deleteRoomIfEmpty(ctx, s, v.GuildID, room.ID)
	}
}

// roomName builds a room name from the lobby's template, falling back to the member's name if it is empty.
func roomName(template string, member *discordgo.Member) string {
	displayName := ""
	if member != nil {
		displayName = member.DisplayName()
	}
	if template == "" {
		return displayName
	}
	return strings.ReplaceAll(template, userPlaceholder, displayName)
}

// deleteRoomIfEmpty deletes channelId if it is a room and nobody is in it.
func deleteRoomIfEmpty(ctx context.Context, s *discordgo.Session, guildId string, channelId string) {
	isRoom, err := model.IsVoiceLobbyRoom(ctx, channelId)
	if err != nil {
		slog.Error("Failed to check for a voice lobby room.", "channelId", channelId, "err", err)
		return
	}
	if !isRoom || hasVoiceMembers(s, guildId, channelId) {
		return
	}

	_, err = s.ChannelDelete(channelId, discordgo.WithContext(ctx))
	var restErr *discordgo.RESTError
	if err != nil && !(errors.As(err, &restErr) && restErr.Message != nil && restErr.Message.Code == discordgo.ErrCodeUnknownChannel) {
		slog.Error("Failed to delete a voice lobby room.", "channelId", channelId, "err", err)
		return
	}

	if err := model.DeleteVoiceLobbyRoom(ctx, channelId); err != nil {
		slog.Error("Failed to delete a voice lobby room record.", "channelId", channelId, "err", err)
	}
}

// hasVoiceMembers reports whether anyone is in channelId according to the state cache.
func hasVoiceMembers(s *discordgo.Session, guildId string, channelId string) bool {
	guild, err := s.State.Guild(guildId)
	if err != nil {
		// Without the cache the room cannot be proven empty, so keep it.
		return true
	}

	s.State.RLock()
	defer s.State.RUnlock()
	for _, voiceState := range guild.VoiceStates {
		if voiceState.ChannelID == channelId {
			return true
		}
	}
	return false
}
