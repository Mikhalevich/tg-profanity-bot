package processor

import "fmt"

func (p *processor) ProcessMessage(msg string, rsp ResponseAction) error {
	mangledMsg := p.replacer.Replace(msg)

	if mangledMsg == msg {
		return nil
	}

	if err := rsp.Edit(mangledMsg); err != nil {
		return fmt.Errorf("msg edit: %w", err)
	}

	return nil
}
