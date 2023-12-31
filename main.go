package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/redis/go-redis/v9"
)

var (
	BOT_TOKEN      = os.Getenv("BUNTY_BOT_TOKEN")
	BOT_ID         = ""
	QUOTE_API      = "https://api.quotable.io/random"
	REDIS_URL      = os.Getenv("REDIS_URL")
	discordSession *discordgo.Session
	redisClient    *redis.Client
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

func activityHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BOT_ID || strings.HasPrefix(m.Content, "!") {
		return
	}

	authorId := m.Author.ID
	authorKey := fmt.Sprintf("AuthorActivityCount::%v", authorId)

	redisClient.Incr(context.Background(), authorKey)
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

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	BOT_ID = s.State.User.ID
	s.StateEnabled = true
}

func messageCountListener(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BOT_ID {
		return
	}

	if strings.HasPrefix(m.Content, "!mc") {
		activityData := redisClient.Get(context.Background(), fmt.Sprintf("AuthorActivityCount::%v", m.Author.ID))
		res, err := activityData.Result()
		if err != nil {
			log.Fatal(err)
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Your recentMessageCount: %v", res))
	}
}

func main() {
	var err error

	opt, err := redis.ParseURL(REDIS_URL)
	if err != nil {
		log.Fatal(err)
	}

	redisClient = redis.NewClient(opt)

	discordSession, err = discordgo.New("Bot " + BOT_TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	discordSession.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsGuildPresences
	discordSession.AddHandler(activityHandler)
	discordSession.AddHandler(pingMessageHandler)
	discordSession.AddHandler(quoteHandler)
	discordSession.AddHandler(mithooHandler)
	discordSession.AddHandler(messageCountListener)
	discordSession.AddHandler(onReady)

	discordSession.Open()
	defer discordSession.Close()

	<-make(chan int)
}
