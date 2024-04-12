package position

import (
	"sort"
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
	sort.Sort(s)
	return s.positions
}

func (s *SortedPositions) Len() int {
	return len(s.positions)
}

func (s *SortedPositions) Swap(i, j int) {
	s.positions[i], s.positions[j] = s.positions[j], s.positions[i]
}

func (s *SortedPositions) Less(i, j int) bool {
	if s.positions[i].Pos < s.positions[j].Pos {
		return true
	}

	if s.positions[i].Pos == s.positions[j].Pos {
		if !s.positions[i].IsEnd || s.positions[j].IsEnd {
			return true
		}

		return false
	}

	return false
}
