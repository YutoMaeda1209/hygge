package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// VoiceLobbyRoom is a voice channel created from a lobby. It is deleted once nobody is in it.
type VoiceLobbyRoom struct {
	Id                  uint `gorm:"primarykey"`
	VoiceLobbyChannelId uint
	VoiceLobbyChannel   VoiceLobbyChannel
	ChannelId           string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// CreateVoiceLobbyRoom records the Discord channel channelId as a room created from the lobby lobbyId.
func CreateVoiceLobbyRoom(ctx context.Context, lobbyId uint, channelId string) error {
	room := VoiceLobbyRoom{VoiceLobbyChannelId: lobbyId, ChannelId: channelId}
	return gorm.G[VoiceLobbyRoom](db).Create(ctx, &room)
}

// IsVoiceLobbyRoom reports whether the Discord channel channelId is a room created from a lobby.
func IsVoiceLobbyRoom(ctx context.Context, channelId string) (bool, error) {
	count, err := gorm.G[VoiceLobbyRoom](db).Where("channel_id = ?", channelId).Count(ctx, "id")
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ListVoiceLobbyRoomChannelIds returns the Discord channel ids of all rooms in the guild guildId.
func ListVoiceLobbyRoomChannelIds(ctx context.Context, guildId string) ([]string, error) {
	var channelIds []string
	err := db.WithContext(ctx).
		Model(&VoiceLobbyRoom{}).
		Joins("JOIN voice_lobby_channels ON voice_lobby_channels.id = voice_lobby_rooms.voice_lobby_channel_id").
		Joins("JOIN servers ON servers.id = voice_lobby_channels.server_id").
		Where("servers.guild_id = ?", guildId).
		Pluck("voice_lobby_rooms.channel_id", &channelIds).Error
	return channelIds, err
}

// DeleteVoiceLobbyRoom deletes the room record for the Discord channel channelId.
// Deleting an unknown channel is not an error.
func DeleteVoiceLobbyRoom(ctx context.Context, channelId string) error {
	_, err := gorm.G[VoiceLobbyRoom](db).Where("channel_id = ?", channelId).Delete(ctx)
	return err
}
