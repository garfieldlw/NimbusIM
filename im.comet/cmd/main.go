package main

import (
	"context"
	"github.com/garfieldlw/NimbusIM/im.comet/grpc/handler"
	"github.com/garfieldlw/NimbusIM/im.comet/internal/service/unique"
	"github.com/garfieldlw/NimbusIM/im.comet/internal/service/ws"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	im_comet "github.com/garfieldlw/NimbusIM/proto/im.comet"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	server := grpc.NewServer()

	im_comet.RegisterImCometServer(server, new(handler.ImCometServer))

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

	id, err := unique.GetServerId(context.Background())
	if err != nil {
		panic("get server id failed: " + err.Error())
	}

	ws.Init(id)
	ws.Run()

	grpcListener, _ := net.Listen("tcp", ":50051")
	if err := server.Serve(grpcListener); err != nil {
		log.Fatal("run service fatal", zap.Error(err))
	}
}
