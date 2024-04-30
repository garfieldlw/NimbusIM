package enum

type PrefixEnum string

const (
	PrefixEnumNone         PrefixEnum = "NONE"
	PrefixEnumConfig       PrefixEnum = "config"
	PrefixEnumUser         PrefixEnum = "user"
	PrefixEnumConversation PrefixEnum = "conversation"
	PrefixEnumMessage      PrefixEnum = "message"
	PrefixEnumRead         PrefixEnum = "read"
)

type ChannelEnum string

const (
	ChannelEnumUser ChannelEnum = "channel_user"
)
