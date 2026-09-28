package bot

import (
	"github.com/YutoMaeda1209/hygge/config"
	"github.com/bwmarrin/discordgo"
)

// Start registers the event handlers and opens a gateway connection with config.Conf.BotToken.
// The caller must close the returned session.
func Start() (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + config.Conf.BotToken)
	if err != nil {
		return nil, err
	}
	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates

	session.AddHandler(onGuildCreate)
	session.AddHandler(onVoiceStateUpdate)
	session.AddHandler(onChannelDelete)

	if err := session.Open(); err != nil {
		return nil, err
	}
	return session, nil
}
