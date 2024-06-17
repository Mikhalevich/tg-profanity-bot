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

type ProfanityStaticSuit struct {
	*suite.Suite
	p *profanity
}

func TestProfanityStaticSuit(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ProfanityStaticSuit{
		Suite: new(suite.Suite),
	})
}

func (s *ProfanityStaticSuit) SetupSuite() {
	words, err := config.BadWords()
	if err != nil {
		s.Fail("get bad words: %v", err)
	}

	s.p = New(matcher.NewAhocorasick(words), replacer.NewStatic("<< censored >>"))
}

func (s *ProfanityStaticSuit) TestAhocorasickStatic() {
	var (
		tests = []struct {
			Msg         string
			ExpectedMsg string
		}{
			{Msg: "hello ass", ExpectedMsg: "hello << censored >>"},
			{Msg: "ass hello", ExpectedMsg: "<< censored >> hello"},
			{Msg: "hello ass from", ExpectedMsg: "hello << censored >> from"},
			{Msg: "ass ass", ExpectedMsg: "<< censored >> << censored >>"},
			{Msg: "ass hello ass from ass", ExpectedMsg: "<< censored >> hello << censored >> from << censored >>"},
			{Msg: "shit", ExpectedMsg: "<< censored >>"},
			{Msg: "shit ass", ExpectedMsg: "<< censored >> << censored >>"},
			{Msg: "ass shit", ExpectedMsg: "<< censored >> << censored >>"},
			{Msg: "ass the shit", ExpectedMsg: "<< censored >> the << censored >>"},
			{Msg: "shit the ass", ExpectedMsg: "<< censored >> the << censored >>"},
			{Msg: "the shit the ass", ExpectedMsg: "the << censored >> the << censored >>"},
			{Msg: "shit the ass the", ExpectedMsg: "<< censored >> the << censored >> the"},
			{Msg: "the shit the ass the", ExpectedMsg: "the << censored >> the << censored >> the"},
			{Msg: "shiasst", ExpectedMsg: "shi<< censored >>t"},
			{Msg: "ashitss", ExpectedMsg: "a<< censored >>ss"},
			{Msg: "ball_sucking", ExpectedMsg: "<< censored >>"},
			{Msg: "HeLLo ErotIC", ExpectedMsg: "HeLLo << censored >>"},
			{Msg: "dick diff cases DICK", ExpectedMsg: "<< censored >> diff cases << censored >>"},
			{
				Msg: strings.Repeat("bDSm test sex WITH fucK diFF boob caSeS ANUS eNd", 10),
				ExpectedMsg: strings.Repeat(
					"<< censored >> test << censored >> WITH << censored >> diFF << censored >> caSeS << censored >> eNd",
					10,
				),
			},
			{Msg: "asssuck", ExpectedMsg: "<< censored >>"},
			{Msg: "assuck", ExpectedMsg: "<< censored >>"},
			{Msg: "no replaces", ExpectedMsg: "no replaces"},
			{Msg: "сискмен", ExpectedMsg: "<< censored >>мен"},
			{Msg: "ребёнок", ExpectedMsg: "ребёнок"},
			{Msg: "тебе", ExpectedMsg: "тебе"},
			{Msg: "себе", ExpectedMsg: "себе"},
			{Msg: "ебет", ExpectedMsg: "<< censored >>"},
			{Msg: "ебёт", ExpectedMsg: "<< censored >>"},
		}
	)

	for _, tc := range tests {
		actual, err := s.p.Replace(context.Background(), "", tc.Msg)
		s.Require().NoError(err)
		s.Require().EqualValues(tc.ExpectedMsg, actual)
	}
}

func initStatic(b *testing.B) *profanity {
	b.Helper()

	words, err := config.BadWords()
	if err != nil {
		b.Fatalf("get bad words: %v", err)
	}

	return New(matcher.NewAhocorasick(words), replacer.NewStatic("<< censored >>"))
}

func BenchmarkAhocorasickStaticPredefined(b *testing.B) {
	var (
		tests = []struct {
			Msg         string
			ExpectedMsg string
		}{
			{Msg: "hello ass", ExpectedMsg: "hello << censored >>"},
			{Msg: "ass hello", ExpectedMsg: "<< censored >> hello"},
			{Msg: "hello ass from", ExpectedMsg: "hello << censored >> from"},
			{Msg: "ass ass", ExpectedMsg: "<< censored >> << censored >>"},
			{Msg: "ass hello ass from ass", ExpectedMsg: "<< censored >> hello << censored >> from << censored >>"},
			{Msg: "shit", ExpectedMsg: "<< censored >>"},
			{Msg: "shit ass", ExpectedMsg: "<< censored >> << censored >>"},
			{Msg: "ass shit", ExpectedMsg: "<< censored >> << censored >>"},
			{Msg: "ass the shit", ExpectedMsg: "<< censored >> the << censored >>"},
			{Msg: "shit the ass", ExpectedMsg: "<< censored >> the << censored >>"},
			{Msg: "the shit the ass", ExpectedMsg: "the << censored >> the << censored >>"},
			{Msg: "shit the ass the", ExpectedMsg: "<< censored >> the << censored >> the"},
			{Msg: "the shit the ass the", ExpectedMsg: "the << censored >> the << censored >> the"},
			{Msg: "shiasst", ExpectedMsg: "shi<< censored >>t"},
			{Msg: "ashitss", ExpectedMsg: "a<< censored >>ss"},
			{Msg: "ball_sucking", ExpectedMsg: "<< censored >>"},
			{Msg: "HeLLo ErotIC", ExpectedMsg: "HeLLo << censored >>"},
			{Msg: "dick diff cases DICK", ExpectedMsg: "<< censored >> diff cases << censored >>"},
			{
				Msg: strings.Repeat("bDSm test sex WITH fucK diFF boob caSeS ANUS eNd", 10),
				ExpectedMsg: strings.Repeat(
					"<< censored >> test << censored >> WITH << censored >> diFF << censored >> caSeS << censored >> eNd",
					10,
				),
			},
			{Msg: "asssuck", ExpectedMsg: "<< censored >>"},
			{Msg: "assuck", ExpectedMsg: "<< censored >>"},
			{Msg: "no replaces", ExpectedMsg: "no replaces"},
			{Msg: "сискмен", ExpectedMsg: "<< censored >>мен"},
			{Msg: "ребёнок", ExpectedMsg: "ребёнок"},
			{Msg: "тебе", ExpectedMsg: "тебе"},
			{Msg: "себе", ExpectedMsg: "себе"},
			{Msg: "ебет", ExpectedMsg: "<< censored >>"},
			{Msg: "ебёт", ExpectedMsg: "<< censored >>"},
		}

		p = initStatic(b)
	)

	for i := 0; i < b.N; i++ {
		for _, tc := range tests {
			if _, err := p.Replace(context.Background(), "", tc.Msg); err != nil {
				b.Fatalf("unexpected error: %v", err)
			}
		}
	}
}

func BenchmarkAhocorasickStaticNoReplacement(b *testing.B) {
	p := initStatic(b)

	for i := 0; i < b.N; i++ {
		if _, err := p.Replace(
			context.Background(),
			"",
			"some text without bad words",
		); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkAhocorasickStaticSmallText(b *testing.B) {
	p := initStatic(b)

	for i := 0; i < b.N; i++ {
		if _, err := p.Replace(
			context.Background(),
			"",
			"some dildo small ass test cock case erotic",
		); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkAhocorasickStaticMediumText(b *testing.B) {
	p := initStatic(b)

	for i := 0; i < b.N; i++ {
		if _, err := p.Replace(
			context.Background(),
			"",
			strings.Repeat("some dildo small ass test cock case erotic", 30),
		); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkAhocorasickStaticLargeText(b *testing.B) {
	p := initStatic(b)

	for i := 0; i < b.N; i++ {
		if _, err := p.Replace(
			context.Background(),
			"",
			strings.Repeat("some dildo small ass test cock case erotic", 30),
		); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
