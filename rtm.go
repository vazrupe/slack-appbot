package appbot

import (
	"github.com/nlopes/slack"
)

type rtmReceiveMessage func(api *slack.Client, msg *slack.MessageEvent)
type rtmReceiveOther func(api *slack.Client, data interface{})

type RTMReceiver struct {
	Filter func(data interface{}) bool

	events chan interface{}

	recvMessage []rtmReceiveMessage
	recvOther   []rtmReceiveOther
}

func newRTMReceiver() *RTMReceiver {
	return &RTMReceiver{
		events: make(chan interface{}),
	}
}

func (r *RTMReceiver) Capacity(size int) *RTMReceiver {
	if size < 0 {
		size = 0
	}
	r.events = make(chan interface{}, size)

	return r
}

func (r *RTMReceiver) run(api *slack.Client) {
	for event := range r.events {
		if r.Filter != nil && !r.Filter(event) {
			continue
		}

		switch ev := event.(type) {
		case *slack.MessageEvent:
			for _, recv := range r.recvMessage {
				go recv(api, ev)
			}

		default:
			for _, recv := range r.recvOther {
				go recv(api, ev)
			}
		}
	}
}

func (r *RTMReceiver) ReceiveMessage(recv ...rtmReceiveMessage) *RTMReceiver {
	if recv != nil {
		r.recvMessage = append(r.recvMessage, recv...)
	}
	return r
}

func (r *RTMReceiver) ReceiveOther(recv ...rtmReceiveOther) *RTMReceiver {
	if recv != nil {
		r.recvOther = append(r.recvOther, recv...)
	}
	return r
}
