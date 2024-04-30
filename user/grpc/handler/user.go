package handler

import (
	"context"
	"github.com/garfieldlw/NimbusIM/proto/common"
	"github.com/garfieldlw/NimbusIM/proto/user"
	user2 "github.com/garfieldlw/NimbusIM/user/internal/service/user"
)

type UserServer struct {
}

func (*UserServer) Upset(ctx context.Context, req *user.UserUpsetRequest) (*user.UserEntity, error) {
	return user2.Upset(ctx, req)
}

func (*UserServer) UpdateStatus(ctx context.Context, req *user.UpdateStatusRequest) (*common.Empty, error) {
	return user2.UpdateStatus(ctx, req)
}
func (*UserServer) DetailById(ctx context.Context, req *common.DetailByIdRequest) (*user.UserEntity, error) {
	return user2.DetailById(ctx, req.Id)
}
