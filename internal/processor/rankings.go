package processor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (p *processor) Rankings(ctx context.Context, info port.MessageInfo, cmdArgs string) error {
	topScores, err := p.rankings.Top(ctx, makeCurrentMonthRankingKey(info.ChatID.String()))
	if err != nil {
		return fmt.Errorf("rankings top: %w", err)
	}

	if err := p.msgSender.Reply(ctx, info, msgFromRankings(topScores)); err != nil {
		return fmt.Errorf("msg reply: %w", err)
	}

	return nil
}

func msgFromRankings(topScores []port.RankingUserScore) string {
	if len(topScores) == 0 {
		return "rankings are empty"
	}

	formattedRankings := make([]string, 0, len(topScores))

	for _, user := range topScores {
		formattedRankings = append(formattedRankings, fmt.Sprintf("%s: %d", user.UserID, user.Score))
	}

	return strings.Join(formattedRankings, "\n")
}

func makeRankingsKey(chatID string, month string) string {
	return fmt.Sprintf("rankings_%s_%s", chatID, month)
}

func makeCurrentMonthRankingKey(chatID string) string {
	return makeRankingsKey(chatID, time.Now().Month().String())
}
