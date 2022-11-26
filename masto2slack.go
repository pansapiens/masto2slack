package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/mattn/go-mastodon"
	"github.com/slack-go/slack"
	config "github.com/spf13/viper"
)

var verbose = true

func getConfig() {
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AddConfigPath("$HOME/.config/masto2slack")
	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func PostStatusToSlack(status *mastodon.Status, webhookURL string) {
	// Statuses that are reblogs ("boosts") have no content, so we need to
	// use the "Reblog" field to get the actual content
	s := status
	//isReblog := false
	if s.Content == "" {
		s = status.Reblog
		//isReblog = true
	}
	if s.Content == "" {
		if verbose {
			fmt.Printf("Status %s has no content or reblog (boost) ?\n", s.ID)
		}
		return
	}

	htmlToMdConverter := md.NewConverter("", true, nil)
	mdMsg, err := htmlToMdConverter.ConvertString(s.Content)
	if err != nil {
		log.Fatal(err)
	}

	attachment := slack.Attachment{
		Text:          mdMsg,
		AuthorName:    s.Account.DisplayName,
		AuthorSubname: s.Account.Username,
		AuthorLink:    s.URL,
		Ts:            json.Number(fmt.Sprint(s.CreatedAt.Unix())),
	}
	slackMsg := &slack.WebhookMessage{Attachments: []slack.Attachment{attachment}}

	err = slack.PostWebhook(webhookURL, slackMsg)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	getConfig()

	slackWebhookURL := config.GetString("slack.webhook_url")

	c := mastodon.NewClient(&mastodon.Config{
		Server: config.GetString("mastodon.server"),
		// ClientID and ClientSecret are not required for reading public statuses
		//ClientID:     config.GetString("mastodon.client_id"),
		//ClientSecret: config.GetString("mastodon.client_secret"),
		AccessToken: config.GetString("mastodon.access_token"),
	})

	lastPostedID := config.GetString("last_status_id")
	if verbose {
		fmt.Printf("last_status_id: %s\n", lastPostedID)
	}

	// Get usr account info
	uacct, err := c.GetAccountCurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Get toots for logged in user, only those since we last checked, earliest first
	pg := mastodon.Pagination{
		SinceID: mastodon.ID(lastPostedID),
	}
	timeline, err := c.GetAccountStatuses(context.Background(), uacct.ID, &pg)
	if err != nil {
		log.Fatal(err)
	}

	for i := len(timeline) - 1; i >= 0; i-- {
		if timeline[i].ID == mastodon.ID(lastPostedID) {
			break
		}
		if verbose {
			fmt.Println("----")
			fmt.Printf("time: %s\n", timeline[i].CreatedAt.String())
			fmt.Printf("ID: %s\n", timeline[i].ID)
			fmt.Printf("webhook: %s\n", slackWebhookURL)
		}

		PostStatusToSlack(timeline[i], slackWebhookURL)
		config.Set("last_status_id", timeline[i].ID)
	}

	// update recorded last_status_id
	config.WriteConfig()
}
