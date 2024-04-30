package unique

import (
	"context"
	"github.com/garfieldlw/NimbusIM/client/unique"
	unique2 "github.com/garfieldlw/NimbusIM/proto/unique"
)

func IdForUser(ctx context.Context) (int64, error) {
	ids, err := unique.Client.GetUniqueId(ctx, 1, unique2.BizType_BizTypeUser)
	if err != nil {
		return 0, err
	}

	return ids[0], nil
}
