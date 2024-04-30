package main

import (
	"github.com/garfieldlw/NimbusIM/chat/grpc/handler"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/garfieldlw/NimbusIM/proto/conversation"
	"github.com/garfieldlw/NimbusIM/proto/message"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	server := grpc.NewServer()

	message.RegisterMessageServer(server, new(handler.MessageServer))

	conversation.RegisterConversationServer(server, new(handler.ConversationServer))

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	go func() {
		for {
			s := <-c
			log.Info("get a signal", zap.String("signal", s.String()))
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
				//cron.Stop()
				server.GracefulStop()
				return
			case syscall.SIGHUP:
			default:
				return
			}
		}
	}()

	grpcListener, _ := net.Listen("tcp", ":50051")
	if err := server.Serve(grpcListener); err != nil {
		log.Fatal("run service fatal", zap.Error(err))
	}
}
