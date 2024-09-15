package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	BOT_TOKEN                    = os.Getenv("BUNTY_BOT_TOKEN")
	BOT_ID                       = ""
	TESTING_SERVER_ID            = "1119670816376889414"
	GUILD_WELCOME_CHANNEL_ID_MAP = map[string]string{"710530782795399230": "711113197054328852"}
	discordSession               *discordgo.Session
)

/*
Message Handlers
----------------------------------------------------------
*/
func pingMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BOT_ID {
		return
	}

	if m.Content == "!ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong !")
	}
}

func mithooHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BOT_ID || m.GuildID != TESTING_SERVER_ID {
		return
	}

	if m.Content == "!mithoo" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Hanji bhai jiiiii")
	}
}

func greetingsHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BOT_ID {
		return
	}

	if m.Content == "!hello" || m.Content == "!hi" {
		s.MessageReactionAdd(m.ChannelID, m.ID, "👋")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hello %v", m.Author.Username))
	}
}

func newJoinerHandler(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	if m.User.ID == BOT_ID {
		return
	}

	welcomeChannelId, exists := GUILD_WELCOME_CHANNEL_ID_MAP[m.GuildID]
	if !exists {
		return
	}

	s.ChannelMessageSend(welcomeChannelId, fmt.Sprintf("Welcome to the server %v", m.User.Username))
}

/*
OnReady
----------------------------------------------------------
*/
func onReady(s *discordgo.Session, r *discordgo.Ready) {
	BOT_ID = s.State.User.ID
	s.StateEnabled = true
}

func main() {
	var err error

	discordSession, err = discordgo.New("Bot " + BOT_TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	discordSession.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsGuildPresences
	discordSession.AddHandler(pingMessageHandler)
	discordSession.AddHandler(mithooHandler)
	discordSession.AddHandler(greetingsHandler)
	discordSession.AddHandler(newJoinerHandler)
	discordSession.AddHandler(onReady)

	discordSession.Open()
	defer discordSession.Close()

	<-make(chan int)
}
