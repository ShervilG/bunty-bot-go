package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	BOT_TOKEN         = os.Getenv("BUNTY_BOT_TOKEN")
	BOT_ID            = ""
	TESTING_SERVER_ID = "1119670816376889414"
	discordSession    *discordgo.Session
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

	if m.Content == "!hello" {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hello %v", m.Author.Username))
	}
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
	discordSession.AddHandler(onReady)

	discordSession.Open()
	defer discordSession.Close()

	<-make(chan int)
}
