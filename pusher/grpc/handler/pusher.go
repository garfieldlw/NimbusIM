package handler

import (
	"context"
	"github.com/garfieldlw/NimbusIM/proto/pusher"
	"github.com/garfieldlw/NimbusIM/pusher/internal/service/kafka"
)

type KafkaPusherServer struct {
}

func (*KafkaPusherServer) PushToKafka(ctx context.Context, req *pusher.PusherKafkaRequest) (*pusher.PusherKafkaResponse, error) {
	return kafka.PushToKafka(ctx, req)
}
func (*KafkaPusherServer) PushToKafkaDelay(ctx context.Context, req *pusher.PusherKafkaDelayRequest) (*pusher.PusherKafkaDelayResponse, error) {
	return kafka.PushToKafkaDelay(ctx, req)
}
