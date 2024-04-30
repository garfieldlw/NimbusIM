package user

import (
	"context"
	"errors"
	"github.com/garfieldlw/NimbusIM/proto/common"
	"github.com/garfieldlw/NimbusIM/proto/user"
	user2 "github.com/garfieldlw/NimbusIM/user/internal/dao/user"
	"github.com/garfieldlw/NimbusIM/user/internal/service/unique"
	"github.com/jinzhu/copier"
)

func Upset(ctx context.Context, req *user.UserUpsetRequest) (*user.UserEntity, error) {
	id, err := unique.IdForUser(ctx)
	if err != nil {
		return nil, err
	}

	err = user2.Upset(ctx, id, req.Email)
	if err != nil {
		return nil, err
	}

	return DetailById(ctx, id)
}

func UpdateStatus(ctx context.Context, req *user.UpdateStatusRequest) (*common.Empty, error) {
	err := user2.UpdateStatus(ctx, req.Id, req.Status)
	if err != nil {
		return nil, err
	}

	_ = user2.AsyncInMemory(ctx, req.Id)

	return new(common.Empty), nil
}

func DetailById(ctx context.Context, id int64) (*user.UserEntity, error) {
	l, err := user2.DetailById(ctx, id)
	if err != nil {
		return nil, err
	}

	if l == nil || l.Id == 0 {
		return nil, errors.New("id is 0")
	}

	item := new(user.UserEntity)
	_ = copier.Copy(item, l)

	return item, nil
}
