package appbot

import (
	"log"

	"github.com/nlopes/slack"
)

type Bot struct {
	Token  string
	Debug  bool
	Logger *log.Logger
}

func New() *Bot {
	return &Bot{}
}

func (b *Bot) Run() {
	var options []slack.Option

	if b.Debug {
		options = append(options, slack.OptionDebug(true))
	}
	if b.Logger != nil {
		options = append(options, slack.OptionLog(b.Logger))
	}

	api := slack.New(b.Token, options...)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	targetModules := modules

	for _, m := range targetModules {
		go m.run(api)
	}

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			log.Printf("Connected: %v\n", ev)

		case *slack.PresenceChangeEvent:
			log.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			log.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			log.Panic(ev)

		case *slack.InvalidAuthEvent:
			log.Println("Invalid credentials")
			return

		default:
			go func() {
				for _, m := range targetModules {
					m.rtm.events <- ev
				}
			}()
		}
	}
}
