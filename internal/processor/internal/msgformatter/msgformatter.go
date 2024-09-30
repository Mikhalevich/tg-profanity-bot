package msgformatter

import (
	"strings"
	"unicode/utf8"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MsgFormatter struct {
	offset      int
	lines       []string
	currentLine []string
	format      []port.Format
}

func New(linesCount int) *MsgFormatter {
	return &MsgFormatter{
		lines:  make([]string, 0, linesCount),
		format: make([]port.Format, 0, linesCount),
	}
}

func (mf *MsgFormatter) AddPlainTextPart(part string) {
	mf.currentLine = append(mf.currentLine, part)
	mf.offset += utf8.RuneCountInString(part)
}

func (mf *MsgFormatter) AddBoldPart(part string) {
	offset := mf.offset
	mf.AddPlainTextPart(part)
	partLen := mf.offset - offset

	mf.format = append(mf.format, port.Format{
		Type:   port.FormatTypeBold,
		Offset: offset,
		Length: partLen,
	})
}

func (mf *MsgFormatter) AddMentionPart(part string, user *tgbotapi.User) {
	offset := mf.offset
	mf.AddPlainTextPart(part)
	partLen := mf.offset - offset

	mf.format = append(mf.format, port.Format{
		Type:   port.FormatTypeMention,
		Offset: offset,
		Length: partLen,
		User:   user,
	})
}

func (mf *MsgFormatter) CompleteLine() {
	mf.currentLine = append(mf.currentLine, "\n")
	mf.offset += 1

	mf.lines = append(mf.lines, strings.Join(mf.currentLine, ""))
	mf.currentLine = mf.currentLine[:0]
}

func (mf *MsgFormatter) ResultString() (string, []port.Format) {
	return strings.Join(mf.lines, ""), mf.format
}
