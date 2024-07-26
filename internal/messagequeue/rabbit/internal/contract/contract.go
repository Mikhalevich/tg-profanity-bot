package contract

type MessageType string

const (
	Message       MessageType = "message"
	CallbackQuery MessageType = "callbackQuery"
)

func (mt MessageType) String() string {
	return string(mt)
}
