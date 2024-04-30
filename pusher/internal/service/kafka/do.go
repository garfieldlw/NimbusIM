package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/garfieldlw/NimbusIM/pkg/delay"
	"github.com/garfieldlw/NimbusIM/proto/pusher"
	"time"
)

func PushToKafka(ctx context.Context, req *pusher.PusherKafkaRequest) (*pusher.PusherKafkaResponse, error) {
	if req.Id < 1 {
		return nil, errors.New("push to kafka, invalid id")
	}

	if len(req.Topic) == 0 {
		return nil, errors.New("push to kafka, invalid topic")
	}

	if len(req.Data) == 0 {
		return nil, errors.New("push to kafka, invalid data")
	}

	service := LoadKafkaService()
	if service == nil {
		return nil, errors.New("push to kafka, load client fail")
	}

	err := service.PushDataToKafka(req.Id, req.Topic, req.Data)
	if err != nil {
		return nil, err
	}

	resp := new(pusher.PusherKafkaResponse)
	resp.Id = req.Id

	return resp, nil
}

func PushToKafkaDelay(ctx context.Context, req *pusher.PusherKafkaDelayRequest) (*pusher.PusherKafkaDelayResponse, error) {
	if req.Id < 1 {
		return nil, errors.New("push to kafka, invalid id")
	}

	if len(req.Topic) == 0 {
		return nil, errors.New("push to kafka, invalid topic")
	}

	if len(req.Data) == 0 {
		return nil, errors.New("push to kafka, invalid data")
	}

	data := new(delay.PayloadDelay)
	data.Id = req.Id
	data.Topic = req.Topic
	data.Payload = req.Data
	data.Push = time.Now().Unix()

	d, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var topic delay.DelayTopic
	switch req.Delay {
	case pusher.DelayEnum_DelayEnum1m:
		{
			topic = delay.DelayTopic1m
		}
	case pusher.DelayEnum_DelayEnum3m:
		{
			topic = delay.DelayTopic3m
		}
	case pusher.DelayEnum_DelayEnum10m:
		{
			topic = delay.DelayTopic10m
		}
	case pusher.DelayEnum_DelayEnum15m:
		{
			topic = delay.DelayTopic15m
		}
	case pusher.DelayEnum_DelayEnum30m:
		{
			topic = delay.DelayTopic30m
		}
	default:
		{
			return nil, errors.New("push to kafka, invalid delay type")
		}
	}
	service := LoadKafkaService()
	if service == nil {
		return nil, errors.New("push to kafka, load client fail")
	}

	err = service.PushDataToKafka(req.Id, string(topic), string(d[:]))
	if err != nil {
		return nil, err
	}

	resp := new(pusher.PusherKafkaDelayResponse)
	resp.Id = req.Id

	return resp, nil
}
