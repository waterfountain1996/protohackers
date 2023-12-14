package skiplist

import "testing"

func TestSkipListInsert(t *testing.T) {
	sl := NewSkipList(4)

	sl.Insert(1, "foo")
	sl.Insert(2, "bar")
	sl.Insert(420, "baz")

	if sl.Head == nil {
		t.Fatal("Expected header not to be nil")
	}

	first := sl.Head.Next[0]
	second := first.Next[0]
	third := second.Next[0]

	if !(first.Score == 1 && first.Value == "foo") {
		t.Fatalf(
			"Expected first.Score == %d (got %d), first.Value == %s (got %s)",
			1,
			first.Score,
			"foo",
			first.Value,
		)
	}

	if !(second.Score == 2 && second.Value == "bar") {
		t.Fatalf(
			"Expected second.Score == %d (got %d), second.Value == %s (got %s)",
			2,
			second.Score,
			"bar",
			second.Value,
		)
	}

	if !(third.Score == 420 && third.Value == "baz") {
		t.Fatalf(
			"Expected third.Score == %d (got %d), third.Value == %s (got %s)",
			420,
			third.Score,
			"baz",
			third.Value,
		)
	}
}

func TestSkipListRangeByScore(t *testing.T) {
	sl := NewSkipList(4)

	sl.Insert(1, "foo")
	sl.Insert(2, "bar")
	sl.Insert(3, "baz")
	sl.Insert(4, "qux")
	sl.Insert(5, "quux")
	sl.Insert(7, "quuux")

	values := sl.RangeByScore(2, 4)
	expected := []string{"bar", "baz", "qux"}
	for idx, value := range values {
		if value.(string) != expected[idx] {
			t.Fatalf("Expected values[%d] == %s, got %s", idx, expected[idx], value)
		}
	}
}
