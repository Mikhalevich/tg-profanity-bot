package processor

func (p *processor) ProcessMessage(msg string, rsp ResponseAction) error {
	mangledMsg := p.replacer.Replace(msg)

	if mangledMsg != msg {
		return rsp.Edit(mangledMsg)
	}

	return nil
}
