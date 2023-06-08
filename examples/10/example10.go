package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack"
	"github.com/yunginnanet/slacker"
)

const (
	errorFormat = "> Custom Error: _%s_"
)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	bot.CustomResponse(NewCustomResponseWriter)

	definition := &slacker.CommandDefinition{
		Description: "Custom!",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("custom")
			response.ReportError(errors.New("oops, an error occurred"))
		},
	}

	bot.Command("custom", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

// NewCustomResponseWriter creates a new ResponseWriter structure
func NewCustomResponseWriter(botCtx slacker.BotContext) slacker.ResponseWriter {
	return &MyCustomResponseWriter{botCtx: botCtx}
}

// MyCustomResponseWriter a custom response writer
type MyCustomResponseWriter struct {
	botCtx slacker.BotContext
}

// ReportError sends back a formatted error message to the channel where we received the event from
func (r *MyCustomResponseWriter) ReportError(err error, options ...slacker.ReportErrorOption) {
	defaults := slacker.NewReportErrorDefaults(options...)

	apiClient := r.botCtx.APIClient()
	event := r.botCtx.Event()

	opts := []slack.MsgOption{
		slack.MsgOptionText(fmt.Sprintf(errorFormat, err.Error()), false),
	}
	if defaults.ThreadResponse {
		opts = append(opts, slack.MsgOptionTS(event.TimeStamp))
	}

	_, _, err = apiClient.PostMessage(event.ChannelID, opts...)
	if err != nil {
		fmt.Printf("failed to report error: %v\n", err)
	}
}

// Reply send a message to the current channel
func (r *MyCustomResponseWriter) Reply(message string, options ...slacker.ReplyOption) error {
	ev := r.botCtx.Event()
	if ev == nil {
		return fmt.Errorf("unable to get message event details")
	}
	return r.Post(ev.ChannelID, message, options...)
}

// Post send a message to a channel
func (r *MyCustomResponseWriter) Post(channel string, message string, options ...slacker.ReplyOption) error {
	defaults := slacker.NewReplyDefaults(options...)

	apiClient := r.botCtx.APIClient()
	ev := r.botCtx.Event()
	if ev == nil {
		return fmt.Errorf("unable to get message event details")
	}

	opts := []slack.MsgOption{
		slack.MsgOptionText(message, false),
		slack.MsgOptionAttachments(defaults.Attachments...),
		slack.MsgOptionBlocks(defaults.Blocks...),
	}

	if defaults.ThreadResponse {
		opts = append(opts, slack.MsgOptionTS(ev.TimeStamp))
	}

	_, _, err := apiClient.PostMessage(
		channel,
		opts...,
	)
	return err
}
