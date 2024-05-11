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

func (p *processor) ProcessMessage(msg string, rsp ResponseAction) error {
	mangledMsg := p.replacer.Replace(msg)

	if mangledMsg != msg {
		return rsp.Edit(mangledMsg)
	}

	return nil
}

func (p *processor) ProcessCommand(cmd string, args string, rsp ResponseAction) error {
	return nil
}
