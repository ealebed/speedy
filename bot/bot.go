package bot

import (
	"errors"
	"fmt"
	"log"
	"strings"

	handlersStore "github.com/ealebed/speedy/handlers"
	"github.com/slack-go/slack"
)

type handler func(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string)

// Run runs Slack Bot
func InitSlack(token string) error {
	if len(token) == 0 {
		log.Printf("Encountered an error: Token was not provided!")
	}

	c := slack.New(token)
	rtm := c.NewRTM()
	err := make(chan error)

	go serveEvents(c, rtm, err)
	go rtm.ManageConnection()
	return <-err
}

func serveEvents(c *slack.Client, rtm *slack.RTM, err chan error) {
	var currentUser string
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				currentUser = fmt.Sprintf("<@%s>", ev.Info.User.ID)
				log.Printf("ConnectedEvent %+v\ncurrentUser:%+v", ev, currentUser)
			case *slack.HelloEvent:
				log.Printf("HelloEvent %+v\n", ev)
			case *slack.InvalidAuthEvent:
				log.Printf("InvalidAuthEvent %+v\n", ev)
				err <- errors.New("auth error")
			case *slack.ConnectionErrorEvent:
				log.Printf("ConnectionErrorEvent %+v\n", ev)
				err <- errors.New("connection error")
			case *slack.MessageEvent:
				log.Printf("MessageEvent %+v\n", ev)
				handleMessageEvent(c, rtm, ev, currentUser)
			}
		}
	}
}

var handlers = map[string]handler{
	"Hi":    handlersStore.GreetHandler,
	"hi":    handlersStore.GreetHandler,
	"Hey":   handlersStore.GreetHandler,
	"hey":   handlersStore.GreetHandler,
	"Hello": handlersStore.GreetHandler,
	"hello": handlersStore.GreetHandler,
}

func handleMessageEvent(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, u string) {
	cmds := strings.Split(ev.Text, " ")
	log.Printf("cmds: %v\n", cmds)
	var key string
	l := len(cmds)
	switch {
	case l == 1:
		key = cmds[0]
	case l >= 2:
		if cmds[0] == u {
			key = cmds[1]
			cmds = cmds[1:]
		} else {
			key = cmds[0]
		}
	default:
		return
	}
	if f, ok := handlers[key]; ok {
		f(c, rtm, ev, cmds)
		return
	}
}
