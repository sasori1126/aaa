package slack

import (
	"axis/ecommerce-backend/configs"
	"github.com/woopla/slack-go-webhook"
)

func NewEvent(config *configs.SlackConfig, author, titleLink, color *string) {
	if author == nil {
		system := "System Monitor"
		author = &system
	}

	if titleLink == nil {
		lnk := "https://axisforestry.com"
		titleLink = &lnk
	}

	if color == nil {
		clr := "#FF0000"
		color = &clr
	}

	attachment1 := slack.Attachment{Color: color, TitleLink: titleLink, Title: &config.Title}
	attachment1.
		AddField(slack.Field{Title: "Author", Value: *author}).
		AddField(slack.Field{Title: "Event", Value: config.Event})

	payload := slack.Payload(config.Message,
		"AxisBot",
		"",
		config.Channel,
		[]slack.Attachment{attachment1})

	slack.Send(config.Webhook, "", payload)
}
