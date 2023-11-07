package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	BOT_TOKEN      = os.Getenv("BUNTY_BOT_TOKEN")
	BOT_ID         = ""
	QUOTE_API      = "https://api.quotable.io/random"
	discordSession *discordgo.Session
)

type Quote struct {
	ID           string   `json:"_id"`
	Content      string   `json:"content"`
	Author       string   `json:"author"`
	Tags         []string `json:"tags"`
	AuthorSlug   string   `json:"authorSlug"`
	Length       int      `json:"length"`
	DateAdded    string   `json:"dateAdded"`
	DateModified string   `json:"dateModified"`
}

func pingMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BOT_ID {
		return
	}

	if m.Content == "!ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong !")
	}
}

func quoteHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BOT_ID {
		return
	}

	if m.Content == "!mot" || m.Content == "!quote" {
		resp, err := http.Get(QUOTE_API)
		if err != nil {
			println(err)
		}

		body, error := io.ReadAll(resp.Body)
		if error != nil {
			println(error)
		}

		var quote Quote
		json.Unmarshal(body, &quote)

		_, _ = s.ChannelMessageSend(m.ChannelID, quote.Content)
	}
}

func mithooHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BOT_ID {
		return
	}

	if m.Content == "!mithoo" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Hanji bhai jiiiii")
	}
}

func whatGameAmIPlayingHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BOT_ID {
		return
	}

	if strings.HasPrefix(m.Content, "!game") {
		presence, err := s.State.Presence(m.GuildID, m.Author.ID)
		if err != nil {
			log.Fatal(err)
		}

		if len(presence.Activities) > 0 {
			for _, activity := range presence.Activities {
				println(activity.ApplicationID, " ", activity.Name)
			}
		}
	}
}

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
	discordSession.AddHandler(quoteHandler)
	discordSession.AddHandler(mithooHandler)
	discordSession.AddHandler(whatGameAmIPlayingHandler)
	discordSession.AddHandler(onReady)

	discordSession.Open()
	defer discordSession.Close()

	<-make(chan int)
}
