package position

import (
	"slices"
)

type Position struct {
	Pos   int
	IsEnd bool
}

type SortedPositions struct {
	positions []*Position
}

func NewSortedPositions() *SortedPositions {
	return &SortedPositions{}
}

func (s *SortedPositions) Append(pos *Position) {
	s.positions = append(s.positions, pos)
}

func (s *SortedPositions) Positions() []*Position {
	slices.SortFunc(s.positions, func(a, b *Position) int {
		if a.Pos < b.Pos {
			return -1
		}

		if a.Pos == b.Pos {
			if !a.IsEnd || b.IsEnd {
				return -1
			}

			return 1
		}

		return 1
	})

	return s.positions
}
