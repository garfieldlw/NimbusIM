package unique

import (
	"context"
	"github.com/garfieldlw/NimbusIM/client/unique"
	unique2 "github.com/garfieldlw/NimbusIM/proto/unique"
)

func GetMessageId(ctx context.Context) (int64, error) {
	ids, err := unique.Client.GetUniqueId(ctx, 1, unique2.BizType_BizTypeMessage)
	if err != nil {
		return 0, err
	}

	return ids[0], nil
}

func GetConversationId(ctx context.Context) (int64, error) {
	ids, err := unique.Client.GetUniqueId(ctx, 1, unique2.BizType_BizTypeConversation)
	if err != nil {
		return 0, err
	}

	return ids[0], nil
}
