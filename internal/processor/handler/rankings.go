package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/msgformatter"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

var (
	monthInfo = []struct {
		Long  string
		Short string
	}{
		{
			Long:  "January",
			Short: "Jan",
		},
		{
			Long:  "February",
			Short: "Feb",
		},
		{
			Long:  "March",
			Short: "Mar",
		},
		{
			Long:  "April",
			Short: "Apr",
		},
		{
			Long:  "May",
			Short: "May",
		},
		{
			Long:  "June",
			Short: "Jun",
		},
		{
			Long:  "July",
			Short: "Jul",
		},
		{
			Long:  "August",
			Short: "Aug",
		},
		{
			Long:  "September",
			Short: "Sep",
		},
		{
			Long:  "October",
			Short: "Oct",
		},
		{
			Long:  "November",
			Short: "Nov",
		},
		{
			Long:  "December",
			Short: "Dec",
		},
	}
)

func (h *handler) Rankings(ctx context.Context, info port.MessageInfo, monthArg string) error {
	var (
		month       = parseMonth(monthArg)
		rankingsKey = makeRankingsKey(info.ChatID.String(), month)
	)

	topScores, err := h.rankings.Top(ctx, rankingsKey)
	if err != nil {
		return fmt.Errorf("rankings top: %w", err)
	}

	msg, format, err := h.makeFormattedRankingsMsg(ctx, month, info.ChatID.Int64(), topScores)
	if err != nil {
		return fmt.Errorf("make ranking msg: %w", err)
	}

	if err := h.msgSender.Reply(ctx, info, msg, convertFormatToOptions(format)...); err != nil {
		return fmt.Errorf("msg reply: %w", err)
	}

	return nil
}

func convertFormatToOptions(format []port.Format) []port.Option {
	if len(format) == 0 {
		return nil
	}

	options := make([]port.Option, 0, len(format))

	for _, f := range format {
		options = append(options, port.WithFormat(f))
	}

	return options
}

func parseMonth(monthArg string) string {
	if monthArg == "" {
		return time.Now().Month().String()
	}

	if month := parseMonthByNumber(monthArg); month != "" {
		return month
	}

	if month := parseMonthByName(monthArg); month != "" {
		return month
	}

	return time.Now().Month().String()
}

func parseMonthByNumber(monthArg string) string {
	monthNum, err := strconv.ParseInt(monthArg, 10, 64)
	if err != nil {
		// skip error
		return ""
	}

	if monthNum < 1 || monthNum > 12 {
		return ""
	}

	return monthInfo[monthNum-1].Long
}

func parseMonthByName(monthArg string) string {
	lowerArg := strings.ToLower(monthArg)
	for _, m := range monthInfo {
		if lowerArg == strings.ToLower(m.Long) ||
			lowerArg == strings.ToLower(m.Short) {
			return m.Long
		}
	}

	return ""
}

func (h *handler) makeFormattedRankingsMsg(
	ctx context.Context,
	month string,
	chatID int64,
	topScores []port.RankingUserScore,
) (string, []port.Format, error) {
	if len(topScores) == 0 {
		return fmt.Sprintf("rankings for %s are empty", month), nil, nil
	}

	formatter := msgformatter.New(len(topScores) + 1)

	formatter.AddBoldPart(month)
	formatter.CompleteLine()

	for i, user := range topScores {
		id, err := port.NewIDFromString(user.UserID)
		if err != nil {
			return "", nil, fmt.Errorf("invalid id %q: %w", user.UserID, err)
		}

		userInfo, err := h.permissionChecker.UserInfo(ctx, chatID, id.Int64())
		if err != nil {
			return "", nil, fmt.Errorf("get user name: %w", err)
		}

		formatter.AddPlainTextPart(fmt.Sprintf("%d: ", i+1))
		formatter.AddMentionPart(userInfo.String(), userInfo)
		formatter.AddPlainTextPart(fmt.Sprintf(": %d", user.Score))
		formatter.CompleteLine()
	}

	msg, format := formatter.ResultString()

	return msg, format, nil
}

func makeRankingsKey(chatID string, month string) string {
	return fmt.Sprintf("rankings:%s_%s", chatID, month)
}

func makeCurrentMonthRankingKey(chatID string) string {
	return makeRankingsKey(chatID, time.Now().Month().String())
}
