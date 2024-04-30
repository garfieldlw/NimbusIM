package message

import (
	"context"
	"github.com/garfieldlw/NimbusIM/chat/internal/dao/conversation"
	"github.com/garfieldlw/NimbusIM/chat/internal/dao/message"
	"github.com/garfieldlw/NimbusIM/pkg/utils"
)

const (
	unreadChanLen = 100000
)

var (
	unreadChan0  chan *unreadInfo
	unreadChan1  chan *unreadInfo
	unreadChan2  chan *unreadInfo
	unreadChan3  chan *unreadInfo
	unreadChan4  chan *unreadInfo
	unreadChan5  chan *unreadInfo
	unreadChan6  chan *unreadInfo
	unreadChan7  chan *unreadInfo
	unreadChan8  chan *unreadInfo
	unreadChan9  chan *unreadInfo
	unreadChan10 chan *unreadInfo
	unreadChan11 chan *unreadInfo
	unreadChan12 chan *unreadInfo
	unreadChan13 chan *unreadInfo
	unreadChan14 chan *unreadInfo
	unreadChan15 chan *unreadInfo
	unreadChan16 chan *unreadInfo
	unreadChan17 chan *unreadInfo
	unreadChan18 chan *unreadInfo
	unreadChan19 chan *unreadInfo
	unreadChan20 chan *unreadInfo
	unreadChan21 chan *unreadInfo
	unreadChan22 chan *unreadInfo
	unreadChan23 chan *unreadInfo
	unreadChan24 chan *unreadInfo
	unreadChan25 chan *unreadInfo
	unreadChan26 chan *unreadInfo
	unreadChan27 chan *unreadInfo
	unreadChan28 chan *unreadInfo
	unreadChan29 chan *unreadInfo
	unreadChan30 chan *unreadInfo
	unreadChan31 chan *unreadInfo
	unreadChan32 chan *unreadInfo
	unreadChan33 chan *unreadInfo
	unreadChan34 chan *unreadInfo
	unreadChan35 chan *unreadInfo
	unreadChan36 chan *unreadInfo
	unreadChan37 chan *unreadInfo
	unreadChan38 chan *unreadInfo
	unreadChan39 chan *unreadInfo
	unreadChan40 chan *unreadInfo

	unreadAtChan0  chan *unreadInfo
	unreadAtChan1  chan *unreadInfo
	unreadAtChan2  chan *unreadInfo
	unreadAtChan3  chan *unreadInfo
	unreadAtChan4  chan *unreadInfo
	unreadAtChan5  chan *unreadInfo
	unreadAtChan6  chan *unreadInfo
	unreadAtChan7  chan *unreadInfo
	unreadAtChan8  chan *unreadInfo
	unreadAtChan9  chan *unreadInfo
	unreadAtChan10 chan *unreadInfo
)

type unreadInfo struct {
	Ctx            context.Context
	MsgId          int64 `json:"msg_id"`
	UserId         int64 `json:"user_id"`
	SendId         int64 `json:"send_id"`
	ConversationId int64 `json:"conversation_id"`
	CreateTime     int64 `json:"create_time"`
}

func Init() {
	unreadChan0 = make(chan *unreadInfo, unreadChanLen)
	unreadChan1 = make(chan *unreadInfo, unreadChanLen)
	unreadChan2 = make(chan *unreadInfo, unreadChanLen)
	unreadChan3 = make(chan *unreadInfo, unreadChanLen)
	unreadChan4 = make(chan *unreadInfo, unreadChanLen)
	unreadChan5 = make(chan *unreadInfo, unreadChanLen)
	unreadChan6 = make(chan *unreadInfo, unreadChanLen)
	unreadChan7 = make(chan *unreadInfo, unreadChanLen)
	unreadChan8 = make(chan *unreadInfo, unreadChanLen)
	unreadChan9 = make(chan *unreadInfo, unreadChanLen)
	unreadChan10 = make(chan *unreadInfo, unreadChanLen)
	unreadChan11 = make(chan *unreadInfo, unreadChanLen)
	unreadChan12 = make(chan *unreadInfo, unreadChanLen)
	unreadChan13 = make(chan *unreadInfo, unreadChanLen)
	unreadChan14 = make(chan *unreadInfo, unreadChanLen)
	unreadChan15 = make(chan *unreadInfo, unreadChanLen)
	unreadChan16 = make(chan *unreadInfo, unreadChanLen)
	unreadChan17 = make(chan *unreadInfo, unreadChanLen)
	unreadChan18 = make(chan *unreadInfo, unreadChanLen)
	unreadChan19 = make(chan *unreadInfo, unreadChanLen)
	unreadChan20 = make(chan *unreadInfo, unreadChanLen)
	unreadChan21 = make(chan *unreadInfo, unreadChanLen)
	unreadChan22 = make(chan *unreadInfo, unreadChanLen)
	unreadChan23 = make(chan *unreadInfo, unreadChanLen)
	unreadChan24 = make(chan *unreadInfo, unreadChanLen)
	unreadChan25 = make(chan *unreadInfo, unreadChanLen)
	unreadChan26 = make(chan *unreadInfo, unreadChanLen)
	unreadChan27 = make(chan *unreadInfo, unreadChanLen)
	unreadChan28 = make(chan *unreadInfo, unreadChanLen)
	unreadChan29 = make(chan *unreadInfo, unreadChanLen)
	unreadChan30 = make(chan *unreadInfo, unreadChanLen)
	unreadChan31 = make(chan *unreadInfo, unreadChanLen)
	unreadChan32 = make(chan *unreadInfo, unreadChanLen)
	unreadChan33 = make(chan *unreadInfo, unreadChanLen)
	unreadChan34 = make(chan *unreadInfo, unreadChanLen)
	unreadChan35 = make(chan *unreadInfo, unreadChanLen)
	unreadChan36 = make(chan *unreadInfo, unreadChanLen)
	unreadChan37 = make(chan *unreadInfo, unreadChanLen)
	unreadChan38 = make(chan *unreadInfo, unreadChanLen)
	unreadChan39 = make(chan *unreadInfo, unreadChanLen)
	unreadChan40 = make(chan *unreadInfo, unreadChanLen)

	unreadAtChan0 = make(chan *unreadInfo, unreadChanLen)
	unreadAtChan1 = make(chan *unreadInfo, unreadChanLen)
	unreadAtChan2 = make(chan *unreadInfo, unreadChanLen)
	unreadAtChan3 = make(chan *unreadInfo, unreadChanLen)
	unreadAtChan4 = make(chan *unreadInfo, unreadChanLen)
	unreadAtChan5 = make(chan *unreadInfo, unreadChanLen)
	unreadAtChan6 = make(chan *unreadInfo, unreadChanLen)
	unreadAtChan7 = make(chan *unreadInfo, unreadChanLen)
	unreadAtChan8 = make(chan *unreadInfo, unreadChanLen)
	unreadAtChan9 = make(chan *unreadInfo, unreadChanLen)
	unreadAtChan10 = make(chan *unreadInfo, unreadChanLen)
}

func updateUnread(ctx context.Context, msgId, userId, sendId, conversationId, createTime int64) {
	c := new(unreadInfo)
	c.Ctx = utils.CopyContext(ctx)
	c.MsgId = msgId
	c.UserId = userId
	c.SendId = sendId
	c.ConversationId = conversationId
	c.CreateTime = createTime

	m := userId % 41

	var uc chan *unreadInfo
	switch m {
	case 0:
		uc = unreadChan0
	case 1:
		uc = unreadChan1
	case 2:
		uc = unreadChan2
	case 3:
		uc = unreadChan3
	case 4:
		uc = unreadChan4
	case 5:
		uc = unreadChan5
	case 6:
		uc = unreadChan6
	case 7:
		uc = unreadChan7
	case 8:
		uc = unreadChan8
	case 9:
		uc = unreadChan9
	case 10:
		uc = unreadChan10
	case 11:
		uc = unreadChan11
	case 12:
		uc = unreadChan12
	case 13:
		uc = unreadChan13
	case 14:
		uc = unreadChan14
	case 15:
		uc = unreadChan15
	case 16:
		uc = unreadChan16
	case 17:
		uc = unreadChan17
	case 18:
		uc = unreadChan18
	case 19:
		uc = unreadChan19
	case 20:
		uc = unreadChan20
	case 21:
		uc = unreadChan21
	case 22:
		uc = unreadChan22
	case 23:
		uc = unreadChan23
	case 24:
		uc = unreadChan24
	case 25:
		uc = unreadChan25
	case 26:
		uc = unreadChan26
	case 27:
		uc = unreadChan27
	case 28:
		uc = unreadChan28
	case 29:
		uc = unreadChan29
	case 30:
		uc = unreadChan30
	case 31:
		uc = unreadChan31
	case 32:
		uc = unreadChan32
	case 33:
		uc = unreadChan33
	case 34:
		uc = unreadChan34
	case 35:
		uc = unreadChan35
	case 36:
		uc = unreadChan36
	case 37:
		uc = unreadChan37
	case 38:
		uc = unreadChan38
	case 39:
		uc = unreadChan39
	case 40:
		uc = unreadChan40
	default:
		uc = unreadChan0
	}

	if len(uc) == unreadChanLen {
		return
	}

	uc <- c
}

func updateUnreadAt(ctx context.Context, msgId, userId, conversationId, createTime int64) {
	c := new(unreadInfo)
	c.Ctx = utils.CopyContext(ctx)
	c.MsgId = msgId
	c.UserId = userId
	c.ConversationId = conversationId
	c.CreateTime = createTime

	m := userId % 11

	var uc chan *unreadInfo
	switch m {
	case 0:
		uc = unreadAtChan0
	case 1:
		uc = unreadAtChan1
	case 2:
		uc = unreadAtChan2
	case 3:
		uc = unreadAtChan3
	case 4:
		uc = unreadAtChan4
	case 5:
		uc = unreadAtChan5
	case 6:
		uc = unreadAtChan6
	case 7:
		uc = unreadAtChan7
	case 8:
		uc = unreadAtChan8
	case 9:
		uc = unreadAtChan9
	case 10:
		uc = unreadAtChan10
	default:
		uc = unreadAtChan0
	}

	if len(uc) == unreadChanLen {
		return
	}

	uc <- c
}

func DoUnread0() {
	for {
		select {
		case cmd := <-unreadChan0:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread1() {
	for {
		select {
		case cmd := <-unreadChan1:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread2() {
	for {
		select {
		case cmd := <-unreadChan2:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread3() {
	for {
		select {
		case cmd := <-unreadChan3:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread4() {
	for {
		select {
		case cmd := <-unreadChan4:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread5() {
	for {
		select {
		case cmd := <-unreadChan5:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread6() {
	for {
		select {
		case cmd := <-unreadChan6:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread7() {
	for {
		select {
		case cmd := <-unreadChan7:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread8() {
	for {
		select {
		case cmd := <-unreadChan8:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread9() {
	for {
		select {
		case cmd := <-unreadChan9:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread10() {
	for {
		select {
		case cmd := <-unreadChan10:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread11() {
	for {
		select {
		case cmd := <-unreadChan11:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread12() {
	for {
		select {
		case cmd := <-unreadChan12:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread13() {
	for {
		select {
		case cmd := <-unreadChan13:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread14() {
	for {
		select {
		case cmd := <-unreadChan14:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread15() {
	for {
		select {
		case cmd := <-unreadChan15:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread16() {
	for {
		select {
		case cmd := <-unreadChan16:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread17() {
	for {
		select {
		case cmd := <-unreadChan17:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread18() {
	for {
		select {
		case cmd := <-unreadChan18:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread19() {
	for {
		select {
		case cmd := <-unreadChan19:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread20() {
	for {
		select {
		case cmd := <-unreadChan20:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread21() {
	for {
		select {
		case cmd := <-unreadChan21:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread22() {
	for {
		select {
		case cmd := <-unreadChan22:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread23() {
	for {
		select {
		case cmd := <-unreadChan23:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread24() {
	for {
		select {
		case cmd := <-unreadChan24:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread25() {
	for {
		select {
		case cmd := <-unreadChan25:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread26() {
	for {
		select {
		case cmd := <-unreadChan26:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread27() {
	for {
		select {
		case cmd := <-unreadChan27:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread28() {
	for {
		select {
		case cmd := <-unreadChan28:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread29() {
	for {
		select {
		case cmd := <-unreadChan29:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread30() {
	for {
		select {
		case cmd := <-unreadChan30:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread31() {
	for {
		select {
		case cmd := <-unreadChan31:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread32() {
	for {
		select {
		case cmd := <-unreadChan32:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread33() {
	for {
		select {
		case cmd := <-unreadChan33:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread34() {
	for {
		select {
		case cmd := <-unreadChan34:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread35() {
	for {
		select {
		case cmd := <-unreadChan35:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread36() {
	for {
		select {
		case cmd := <-unreadChan36:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread37() {
	for {
		select {
		case cmd := <-unreadChan37:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread38() {
	for {
		select {
		case cmd := <-unreadChan38:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread39() {
	for {
		select {
		case cmd := <-unreadChan39:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnread40() {
	for {
		select {
		case cmd := <-unreadChan40:
			_ = updateUnreadInfo(cmd)
		}
	}
}

func DoUnreadAt0() {
	for {
		select {
		case cmd := <-unreadAtChan0:
			_ = message.UnreadAt(cmd.Ctx, cmd.MsgId, cmd.UserId, cmd.ConversationId, cmd.CreateTime)
		}
	}
}

func DoUnreadAt1() {
	for {
		select {
		case cmd := <-unreadAtChan1:
			_ = message.UnreadAt(cmd.Ctx, cmd.MsgId, cmd.UserId, cmd.ConversationId, cmd.CreateTime)
		}
	}
}

func DoUnreadAt2() {
	for {
		select {
		case cmd := <-unreadAtChan2:
			_ = message.UnreadAt(cmd.Ctx, cmd.MsgId, cmd.UserId, cmd.ConversationId, cmd.CreateTime)
		}
	}
}

func DoUnreadAt3() {
	for {
		select {
		case cmd := <-unreadAtChan3:
			_ = message.UnreadAt(cmd.Ctx, cmd.MsgId, cmd.UserId, cmd.ConversationId, cmd.CreateTime)
		}
	}
}

func DoUnreadAt4() {
	for {
		select {
		case cmd := <-unreadAtChan4:
			_ = message.UnreadAt(cmd.Ctx, cmd.MsgId, cmd.UserId, cmd.ConversationId, cmd.CreateTime)
		}
	}
}

func DoUnreadAt5() {
	for {
		select {
		case cmd := <-unreadAtChan5:
			_ = message.UnreadAt(cmd.Ctx, cmd.MsgId, cmd.UserId, cmd.ConversationId, cmd.CreateTime)
		}
	}
}

func DoUnreadAt6() {
	for {
		select {
		case cmd := <-unreadAtChan6:
			_ = message.UnreadAt(cmd.Ctx, cmd.MsgId, cmd.UserId, cmd.ConversationId, cmd.CreateTime)
		}
	}
}

func DoUnreadAt7() {
	for {
		select {
		case cmd := <-unreadAtChan7:
			_ = message.UnreadAt(cmd.Ctx, cmd.MsgId, cmd.UserId, cmd.ConversationId, cmd.CreateTime)
		}
	}
}

func DoUnreadAt8() {
	for {
		select {
		case cmd := <-unreadAtChan8:
			_ = message.UnreadAt(cmd.Ctx, cmd.MsgId, cmd.UserId, cmd.ConversationId, cmd.CreateTime)
		}
	}
}

func DoUnreadAt9() {
	for {
		select {
		case cmd := <-unreadAtChan9:
			_ = message.UnreadAt(cmd.Ctx, cmd.MsgId, cmd.UserId, cmd.ConversationId, cmd.CreateTime)
		}
	}
}

func DoUnreadAt10() {
	for {
		select {
		case cmd := <-unreadAtChan10:
			_ = message.UnreadAt(cmd.Ctx, cmd.MsgId, cmd.UserId, cmd.ConversationId, cmd.CreateTime)
		}
	}
}

func updateUnreadInfo(cmd *unreadInfo) error {
	_ = conversation.UpdateConversationList(cmd.Ctx, cmd.ConversationId, cmd.UserId, cmd.CreateTime)
	if cmd.UserId == cmd.SendId {
		return nil
	}

	_ = message.Unread(cmd.Ctx, cmd.MsgId, cmd.UserId, cmd.ConversationId, cmd.CreateTime)
	return nil
}
