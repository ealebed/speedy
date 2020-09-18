package main

import (
	"os"

	"github.com/ealebed/speedy/bot"
)

func main() {
	var slackToken string = os.Getenv("SLACK_TOKEN")
	bot.InitSlack(slackToken)
}
