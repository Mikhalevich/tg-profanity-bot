package replacer

import (
	"strings"
	"unicode/utf8"
)

type dynamicSymbol struct {
	symbol byte
}

func NewDynamicSymbol(symbol byte) *dynamicSymbol {
	return &dynamicSymbol{
		symbol: symbol,
	}
}

func (ds *dynamicSymbol) Replace(text string) string {
	return strings.Repeat(string(ds.symbol), utf8.RuneCountInString(text))
}
