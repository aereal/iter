package seq2_test

import (
	"iter"
	"reflect"
	"slices"
	"testing"

	"github.com/aereal/iter/seq2"
)

func TestFlip(t *testing.T) {
	testCases := []struct {
		name  string
		input []pair[string, int]
		want  []pair[int, string]
	}{
		{
			name:  "values",
			input: []pair[string, int]{{"a", 1}, {"b", 2}, {"c", 3}},
			want:  []pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}},
		},
		{
			name:  "empty",
			input: nil,
			want:  nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := seq2.Flip(iteratePairs(tc.input))
			gotSeq := slices.Collect(collectPairs(got))
			if !reflect.DeepEqual(gotSeq, tc.want) {
				t.Errorf("value mismatch\n\twant: %#v\n\t got: %#v", tc.want, got)
			}
		})
	}
}

type pair[A, B any] struct {
	A A
	B B
}

func collectPairs[A, B any](it iter.Seq2[A, B]) iter.Seq[pair[A, B]] {
	return func(yield func(pair[A, B]) bool) {
		for a, b := range it {
			if !yield(pair[A, B]{a, b}) {
				return
			}
		}
	}
}

func iteratePairs[A, B any](pairs []pair[A, B]) iter.Seq2[A, B] {
	return func(yield func(A, B) bool) {
		for _, pair := range pairs {
			if !yield(pair.A, pair.B) {
				return
			}
		}
	}
}
