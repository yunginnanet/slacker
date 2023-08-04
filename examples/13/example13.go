package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack/socketmode"

	"github.com/yunginnanet/slacker"
)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	bot.Init(func() {
		log.Println("Connected!")
	})

	bot.Err(func(err error) {
		log.Println(err)
	})

	bot.DefaultCommand(func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		response.Reply("Say what?")
	})

	bot.DefaultEvent(func(event interface{}) {
		fmt.Println(event)
	})

	bot.DefaultInnerEvent(func(ctx context.Context, evt interface{}, request *socketmode.Request) {
		fmt.Printf("Handling inner event: %s", evt)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
