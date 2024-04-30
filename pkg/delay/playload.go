package delay

// PayloadDelay 延迟队列数据
type PayloadDelay struct {
	Push    int64  `json:"push"`    // 秒
	Payload string `json:"payload"` //
	Topic   string `json:"topic"`   // 到期之后放入的队列
	Id      int64  `json:"id"`
}
