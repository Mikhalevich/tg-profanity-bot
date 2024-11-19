package testmatcher

import (
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

type buttonMatcher struct {
	caption string
}

func NewButtonMatcher(caption string) buttonMatcher {
	return buttonMatcher{
		caption: caption,
	}
}

func (m buttonMatcher) Matches(x interface{}) bool {
	argButton, ok := x.(*port.Button)
	if !ok {
		return false
	}

	if m.caption != argButton.Caption {
		return false
	}

	return true
}

func (m buttonMatcher) String() string {
	return fmt.Sprintf("matches button by caption %v", m.caption)
}
