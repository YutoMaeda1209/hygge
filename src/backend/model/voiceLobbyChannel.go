package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// VoiceLobbyChannel is a voice channel that, when joined, creates a new voice channel named
// NewChannelName and moves the joining user into it.
type VoiceLobbyChannel struct {
	Id             uint `gorm:"primarykey"`
	ServerId       uint
	Server         Server
	ChannelId      string
	NewChannelName string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// FindByChannelId loads the lobby for the Discord channel channelId into lobby.
// It returns gorm.ErrRecordNotFound if the channel is not a lobby.
func (lobby *VoiceLobbyChannel) FindByChannelId(ctx context.Context, channelId string) error {
	findLobby, err := gorm.G[VoiceLobbyChannel](db).Where("channel_id = ?", channelId).First(ctx)
	if err != nil {
		return err
	}
	*lobby = findLobby
	return nil
}
