package replacer

import (
	"strings"
	"unicode/utf8"
)

type dynamic struct {
	wildcard string
}

func NewDynamic(wildcard string) *dynamic {
	return &dynamic{
		wildcard: wildcard,
	}
}

func (d *dynamic) Replace(text string) string {
	return strings.Repeat(d.wildcard, utf8.RuneCountInString(text))
}
