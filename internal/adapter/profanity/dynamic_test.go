package profanity

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/profanity/matcher"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/profanity/replacer"
	"github.com/Mikhalevich/tg-profanity-bot/internal/config"
)

type ProfanityDynamicSuit struct {
	*suite.Suite
	p *profanity
}

func TestProfanityDynamicSuit(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ProfanityDynamicSuit{
		Suite: new(suite.Suite),
	})
}

func (s *ProfanityDynamicSuit) SetupSuite() {
	words, err := config.BadWords()
	if err != nil {
		s.Fail("get bad words: %v", err)
	}

	s.p = New(matcher.NewAhocorasick(words), replacer.NewDynamic("*"))
}

func (s *ProfanityDynamicSuit) TestAhocorasickDynamic() {
	var (
		tests = []struct {
			Msg         string
			ExpectedMsg string
		}{
			{Msg: "hello ass", ExpectedMsg: "hello ***"},
			{Msg: "ass hello", ExpectedMsg: "*** hello"},
			{Msg: "hello ass from", ExpectedMsg: "hello *** from"},
			{Msg: "ass ass", ExpectedMsg: "*** ***"},
			{Msg: "ass hello ass from ass", ExpectedMsg: "*** hello *** from ***"},
			{Msg: "shit", ExpectedMsg: "****"},
			{Msg: "shit ass", ExpectedMsg: "**** ***"},
			{Msg: "ass shit", ExpectedMsg: "*** ****"},
			{Msg: "ass the shit", ExpectedMsg: "*** the ****"},
			{Msg: "shit the ass", ExpectedMsg: "**** the ***"},
			{Msg: "the shit the ass", ExpectedMsg: "the **** the ***"},
			{Msg: "shit the ass the", ExpectedMsg: "**** the *** the"},
			{Msg: "the shit the ass the", ExpectedMsg: "the **** the *** the"},
			{Msg: "shiasst", ExpectedMsg: "shi***t"},
			{Msg: "ashitss", ExpectedMsg: "a****ss"},
			{Msg: "ball_sucking", ExpectedMsg: "************"},
			{Msg: "HeLLo ErotIC", ExpectedMsg: "HeLLo ******"},
			{Msg: "dick diff cases DICK", ExpectedMsg: "**** diff cases ****"},
			{
				Msg:         strings.Repeat("bDSm test sex WITH fucK diFF boob caSeS ANUS eNd", 10),
				ExpectedMsg: strings.Repeat("**** test *** WITH **** diFF **** caSeS **** eNd", 10),
			},
			{Msg: "asssuck", ExpectedMsg: "*******"},
			{Msg: "assuck", ExpectedMsg: "******"},
			{Msg: "no replaces", ExpectedMsg: "no replaces"},
			{Msg: "сискмен", ExpectedMsg: "****мен"},
			{Msg: "ребёнок", ExpectedMsg: "ребёнок"},
			{Msg: "тебе", ExpectedMsg: "тебе"},
			{Msg: "себе", ExpectedMsg: "себе"},
			{Msg: "ебет", ExpectedMsg: "****"},
			{Msg: "ебёт", ExpectedMsg: "****"},
		}
	)

	for _, tc := range tests {
		actual := s.p.Replace(context.Background(), tc.Msg)
		s.Require().EqualValues(tc.ExpectedMsg, actual)
	}
}

func initDynamic(b *testing.B) *profanity {
	b.Helper()

	words, err := config.BadWords()
	if err != nil {
		b.Fatalf("get bad words: %v", err)
	}

	return New(matcher.NewAhocorasick(words), replacer.NewDynamic("*"))
}

func BenchmarkAhocorasickDynamicPredefined(b *testing.B) {
	var (
		tests = []struct {
			Msg         string
			ExpectedMsg string
		}{
			{Msg: "hello ass", ExpectedMsg: "hello ***"},
			{Msg: "ass hello", ExpectedMsg: "*** hello"},
			{Msg: "hello ass from", ExpectedMsg: "hello *** from"},
			{Msg: "ass ass", ExpectedMsg: "*** ***"},
			{Msg: "ass hello ass from ass", ExpectedMsg: "*** hello *** from ***"},
			{Msg: "shit", ExpectedMsg: "****"},
			{Msg: "shit ass", ExpectedMsg: "**** ***"},
			{Msg: "ass shit", ExpectedMsg: "*** ****"},
			{Msg: "ass the shit", ExpectedMsg: "*** the ****"},
			{Msg: "shit the ass", ExpectedMsg: "**** the ***"},
			{Msg: "the shit the ass", ExpectedMsg: "the **** the ***"},
			{Msg: "shit the ass the", ExpectedMsg: "**** the *** the"},
			{Msg: "the shit the ass the", ExpectedMsg: "the **** the *** the"},
			{Msg: "shiasst", ExpectedMsg: "shi***t"},
			{Msg: "ashitss", ExpectedMsg: "a****ss"},
			{Msg: "ball_sucking", ExpectedMsg: "************"},
			{Msg: "HeLLo ErotIC", ExpectedMsg: "HeLLo ******"},
			{Msg: "dick diff cases DICK", ExpectedMsg: "**** diff cases ****"},
			{
				Msg:         strings.Repeat("bDSm test sex WITH fucK diFF boob caSeS ANUS eNd", 10),
				ExpectedMsg: strings.Repeat("**** test *** WITH **** diFF **** caSeS **** eNd", 10),
			},
			{Msg: "asssuck", ExpectedMsg: "*******"},
			{Msg: "assuck", ExpectedMsg: "******"},
			{Msg: "no replaces", ExpectedMsg: "no replaces"},
			{Msg: "сискмен", ExpectedMsg: "****мен"},
			{Msg: "ребёнок", ExpectedMsg: "ребёнок"},
			{Msg: "тебе", ExpectedMsg: "тебе"},
			{Msg: "себе", ExpectedMsg: "себе"},
			{Msg: "ебет", ExpectedMsg: "****"},
			{Msg: "ебёт", ExpectedMsg: "****"},
		}

		p = initDynamic(b)
	)

	for i := 0; i < b.N; i++ {
		for _, tc := range tests {
			p.Replace(context.Background(), tc.Msg)
		}
	}
}

func BenchmarkAhocorasickDynamicNoReplacement(b *testing.B) {
	p := initDynamic(b)

	for i := 0; i < b.N; i++ {
		p.Replace(context.Background(), "some text without bad words")
	}
}

func BenchmarkAhocorasickDynamicSmallText(b *testing.B) {
	p := initDynamic(b)

	for i := 0; i < b.N; i++ {
		p.Replace(context.Background(), "some dildo small ass test cock case erotic")
	}
}

func BenchmarkAhocorasickDynamicMediumText(b *testing.B) {
	p := initDynamic(b)

	for i := 0; i < b.N; i++ {
		p.Replace(context.Background(), strings.Repeat("some dildo small ass test cock case erotic", 30))
	}
}

func BenchmarkAhocorasickDynamicLargeText(b *testing.B) {
	p := initDynamic(b)

	for i := 0; i < b.N; i++ {
		p.Replace(context.Background(), strings.Repeat("some dildo small ass test cock case erotic", 30))
	}
}
