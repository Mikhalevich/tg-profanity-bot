package processor

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

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

func (p *processor) Rankings(ctx context.Context, info port.MessageInfo, monthArg string) error {
	var (
		month       = parseMonth(monthArg)
		rankingsKey = makeRankingsKey(info.ChatID.String(), month)
	)

	topScores, err := p.rankings.Top(ctx, rankingsKey)
	if err != nil {
		return fmt.Errorf("rankings top: %w", err)
	}

	msg, err := p.makeRankingsMsg(ctx, month, info.ChatID.Int64(), topScores)
	if err != nil {
		return fmt.Errorf("make ranking msg: %w", err)
	}

	if err := p.msgSender.Reply(ctx, info, msg); err != nil {
		return fmt.Errorf("msg reply: %w", err)
	}

	return nil
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

func (p *processor) makeRankingsMsg(
	ctx context.Context,
	month string,
	chatID int64,
	topScores []port.RankingUserScore,
) (string, error) {
	if len(topScores) == 0 {
		return fmt.Sprintf("rankings for %s are empty", month), nil
	}

	formattedRankings := make([]string, 0, len(topScores)+1)

	formattedRankings = append(formattedRankings, month)

	for i, user := range topScores {
		id, err := port.NewIDFromString(user.UserID)
		if err != nil {
			return "", fmt.Errorf("invalid id %q: %w", user.UserID, err)
		}

		userName, err := p.permissionChecker.UserName(ctx, chatID, id.Int64())
		if err != nil {
			return "", fmt.Errorf("get user name: %w", err)
		}

		formattedRankings = append(formattedRankings, fmt.Sprintf("%d: %s: %d", i+1, userName, user.Score))
	}

	return strings.Join(formattedRankings, "\n"), nil
}

func makeRankingsKey(chatID string, month string) string {
	return fmt.Sprintf("rankings:%s_%s", chatID, month)
}

func makeCurrentMonthRankingKey(chatID string) string {
	return makeRankingsKey(chatID, time.Now().Month().String())
}
