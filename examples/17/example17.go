package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yunginnanet/slacker"
)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	bot.SanitizeEventText(func(text string) string {
		fmt.Println("My slack bot does not like backticks!")
		return strings.ReplaceAll(text, "`", "")
	})

	bot.Command("my-command", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("it works!")
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
