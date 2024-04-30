package main

import (
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/garfieldlw/NimbusIM/proto/unique"
	"github.com/garfieldlw/NimbusIM/unique/grpc/handler"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	server := grpc.NewServer()

	unique.RegisterUniqueServer(server, new(handler.UniqueServer))

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

	grpcListener, _ := net.Listen("tcp", ":50051")
	if err := server.Serve(grpcListener); err != nil {
		log.Fatal("run service fatal", zap.Error(err))
	}
}
