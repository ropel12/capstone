package pkg

import (
	"github.com/education-hub/BE/config"
	"github.com/education-hub/BE/errorr"
	"github.com/nsqio/go-nsq"
)

type NSQProducer struct {
	Producer *nsq.Producer
	Env      config.NSQConfig
}

func (np *NSQProducer) Publish(Topic string, message []byte) error {
	switch Topic {
	case "1":
		return np.Producer.Publish(np.Env.Topic, message)
	case "2":
		return np.Producer.Publish(np.Env.Topic2, message)
	case "3":
		return np.Producer.Publish(np.Env.Topic3, message)
	case "4":
		return np.Producer.Publish(np.Env.Topic4, message)
	case "5":
		return np.Producer.Publish(np.Env.Topic5, message)
	case "6":
		return np.Producer.Publish(np.Env.Topic6, message)
	case "7":
		return np.Producer.Publish(np.Env.Topic7, message)
	case "8":
		return np.Producer.Publish(np.Env.Topic8, message)
	case "9":
		return np.Producer.Publish(np.Env.Topic9, message)
	}
	return errorr.NewBad("Topic not available")
}

func (np *NSQProducer) Stop() {
	np.Producer.Stop()
}
