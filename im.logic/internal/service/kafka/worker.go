package kafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	kafka_cli "github.com/garfieldlw/NimbusIM/im.logic/internal/model/dao/kafka-cli"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"go.uber.org/zap"
)

type KafkaConsumerConfig struct {
	Group string `json:"group"`
	Topic string `json:"topic"`
}

type HandlerConfig struct {
	Name      string
	CloseChan chan *struct{}
}

var closeChan *HandlerConfig
var consumers *kafka.Consumer

func RunSub() error {
	kafkaConfig := &KafkaConsumerConfig{
		Group: "group-topic-im-msg",
		Topic: "topic-im-msg",
	}

	if closeChan != nil {
		return nil
	}

	cc := new(HandlerConfig)
	cc.Name = kafkaConfig.Group
	cc.CloseChan = make(chan *struct{})

	go runHandler(kafkaConfig, cc.CloseChan)

	closeChan = cc

	return nil
}

func Close() {
	log.Info("close consumer")

	if consumers == nil {
		return
	}

	err := consumers.Close()
	if err != nil {
		log.Error("close consumer error", zap.Error(err))
	}

	close(closeChan.CloseChan)
}

func runHandler(c *KafkaConsumerConfig, closeChan chan *struct{}) {
	errChan := make(chan error)
	go cron(c, closeChan, errChan)

	for {
		select {
		case workErr := <-errChan:
			log.Warn("worker routine stopped", zap.Error(workErr))
			go cron(c, closeChan, errChan)
		case <-closeChan:
			return
		}
	}
}

func cron(c *KafkaConsumerConfig, closeChan chan *struct{}, errChan chan error) {
	if consumers != nil {
		_ = consumers.Close()
	}

	client, err := kafka_cli.LoadKafkaClientConsumerGroup(c.Group)
	log.Debug("init consumer group", zap.Any("config", *c), zap.Error(err))
	if err != nil {
		errChan <- err
		return
	}

	consumers = client

	err = client.SubscribeTopics([]string{c.Topic}, nil)
	if err != nil {
		errChan <- err
		return
	}

	for {
		select {
		case <-closeChan:
			return

		default:
			ev := client.Poll(1000)
			if ev == nil {
				continue
			}

			switch msg := ev.(type) {
			case *kafka.Message:
				// Process the message received.
				log.Info("consumer group msg", zap.String("offset", msg.TopicPartition.Offset.String()), zap.Int32("partition", msg.TopicPartition.Partition), zap.String("key", string(msg.Key)), zap.String("value", string(msg.Value)))

				err = Do(context.Background(), *msg.TopicPartition.Topic, msg.Key, msg.Value)
				if err == nil {
					_, err = client.CommitMessage(msg)
					if err != nil {
						errChan <- err
						return
					}
				} else {
					log.Info("kafka message error", zap.String("offset", msg.TopicPartition.Offset.String()), zap.Int32("partition", msg.TopicPartition.Partition), zap.String("key", string(msg.Key)), zap.String("value", string(msg.Value)), zap.Error(err))
					_, err = client.CommitOffsets([]kafka.TopicPartition{msg.TopicPartition})
					if err != nil {
						errChan <- err
						return
					}

					_, err = client.SeekPartitions([]kafka.TopicPartition{msg.TopicPartition})
					if err != nil {
						errChan <- err
						return
					}
				}

			case kafka.Error:
				errChan <- msg
				return
			default:
				log.Warn("kafka message error", zap.Any("value", msg))
			}
		}
	}
}
