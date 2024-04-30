package main

import (
	"context"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/garfieldlw/NimbusIM/proto/user"
	"github.com/garfieldlw/NimbusIM/user/grpc/handler"
	"github.com/garfieldlw/NimbusIM/user/internal/service/channel"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	server := grpc.NewServer()

	user.RegisterUserServer(server, new(handler.UserServer))

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	go func() {
		for {
			s := <-c
			log.Info("get a signal", zap.String("signal", s.String()))
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
				server.GracefulStop()
				return
			case syscall.SIGHUP:
			default:
				return
			}
		}
	}()

	ctx := context.Background()
	go channel.RedisSubscribeUser(ctx)

	grpcListener, _ := net.Listen("tcp", ":50051")
	if err := server.Serve(grpcListener); err != nil {
		log.Fatal("run service fatal", zap.Error(err))
	}
}
