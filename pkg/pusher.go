package pkg

import (
	"github.com/education-hub/BE/config"
	"github.com/education-hub/BE/errorr"
	"github.com/pusher/pusher-http-go"
)

type Pusher struct {
	Client *pusher.Client
	Env    config.PusherConfig
}

func (p *Pusher) Publish(data any, event int) error {
	switch event {
	case 1:
		return p.Client.Trigger(p.Env.Channel, p.Env.Event1, data)
	case 2:
		return p.Client.Trigger(p.Env.Channel, p.Env.Event2, data)
	}
	return errorr.NewBad("Event Not Available")
}
