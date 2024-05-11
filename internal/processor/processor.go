package processor

type TextReplacer interface {
	Replace(msg string) string
}

type ResponseAction interface {
	Send(msg string) error
	Edit(msg string) error
}

type processor struct {
	replacer TextReplacer
}

func New(replacer TextReplacer) *processor {
	return &processor{
		replacer: replacer,
	}
}
