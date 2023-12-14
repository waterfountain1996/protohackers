package skiplist

import "math/rand"

type Node struct {
	Score int
	Value interface{}
	Next  []*Node
}

func newNode(score int, value interface{}, height int) *Node {
	return &Node{
		Score: score,
		Value: value,
		Next:  make([]*Node, height),
	}
}

type SkipList struct {
	Head              *Node
	MaxHeight, Height int
}

func NewSkipList(maxHeight int) *SkipList {
	return &SkipList{
		Head:      newNode(0, nil, maxHeight),
		MaxHeight: maxHeight,
		Height:    1,
	}
}

func (sl *SkipList) randLevel() int {
	level := 1
	for rand.Intn(2) == 0 && level < sl.MaxHeight {
		level++
	}
	return level
}

func (sl *SkipList) Insert(score int, value interface{}) {
	tower := make([]*Node, sl.MaxHeight)
	node := sl.Head

	for level := sl.Height - 1; level >= 0; level-- {
		for node.Next[level] != nil && node.Next[level].Score <= score {
			node = node.Next[level]
		}
		tower[level] = node
	}

	newHeight := sl.randLevel()

	if newHeight > sl.Height {
		for level := sl.Height; level < newHeight; level++ {
			tower[level] = sl.Head
		}
		sl.Height = newHeight
	}

	toInsert := newNode(score, value, newHeight)

	for level := 0; level < newHeight; level++ {
		toInsert.Next[level] = tower[level].Next[level]
		tower[level].Next[level] = toInsert
	}
}

func (sl *SkipList) RangeByScore(mn, mx int) []interface{} {
	values := []interface{}{}
	current := sl.Head

	for level := sl.Height - 1; level >= 0; level-- {
		for current.Next[level] != nil && current.Next[level].Score < mn {
			current = current.Next[level]
		}
	}

	current = current.Next[0]

	for current != nil && current.Score <= mx {
		values = append(values, current.Value)
		current = current.Next[0]
	}

	return values
}
