package seq_test

import (
	"iter"
	"reflect"
	"testing"

	"github.com/aereal/iter/seq"
)

func TestDrop(t *testing.T) {
	testCases := []struct {
		name     string
		input    iter.Seq[int]
		inputNum int
		want     []int
	}{
		{name: "ok", input: list(1, 2, 3), inputNum: 2, want: []int{3}},
		{name: "The N is greater than input iterator's length", input: list(1, 2, 3), inputNum: 3, want: []int{}},
		{name: "The N is negative", input: list(1, 2, 3), inputNum: -1, want: []int{1, 2, 3}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := seq.Drop(tc.input, tc.inputNum)
			if gv := values(got); !reflect.DeepEqual(gv, tc.want) {
				t.Errorf("result mismatch:\n\twant: %#v\n\t got: %#v", tc.want, gv)
			}
		})
	}
}

func list[T any](xs ...T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, x := range xs {
			if !yield(x) {
				break
			}
		}
	}
}

func values[T any](s iter.Seq[T]) []T {
	ret := make([]T, 0)
	for v := range s {
		ret = append(ret, v)
	}
	return ret
}
