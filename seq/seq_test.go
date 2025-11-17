package seq_test

import (
	"fmt"
	"iter"
	"reflect"
	"slices"
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

func TestDropWhile(t *testing.T) {
	testCases := []struct {
		name  string
		pred  func(int) bool
		input iter.Seq[int]
		want  []int
	}{
		{name: "ok", input: list(1, 2, 3), pred: func(i int) bool { return i <= 2 }, want: []int{3}},
		{name: "true -> false -> true", input: list(1, 2, 3), pred: func(i int) bool { return i%2 != 0 }, want: []int{2, 3}},
		{name: "all elements ignored", input: list(1, 2, 3), pred: func(i int) bool { return i < 10 }, want: []int{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := seq.DropWhile(tc.input, tc.pred)
			if gv := values(got); !reflect.DeepEqual(gv, tc.want) {
				t.Errorf("result mismatch:\n\twant: %#v\n\t got: %#v", tc.want, gv)
			}
		})
	}
}

func TestTake(t *testing.T) {
	testCases := []struct {
		name     string
		input    iter.Seq[int]
		inputNum int
		want     []int
	}{
		{name: "ok", input: list(1, 2, 3), inputNum: 2, want: []int{1, 2}},
		{name: "The N is greater than input iterator's length", input: list(1, 2, 3), inputNum: 4, want: []int{1, 2, 3}},
		{name: "The N is negative", input: list(1, 2, 3), inputNum: -1, want: []int{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := seq.Take(tc.input, tc.inputNum)
			if gv := values(got); !reflect.DeepEqual(gv, tc.want) {
				t.Errorf("result mismatch:\n\twant: %#v\n\t got: %#v", tc.want, gv)
			}
		})
	}
}

func TestTakeWhile(t *testing.T) {
	testCases := []struct {
		name  string
		pred  func(int) bool
		input iter.Seq[int]
		want  []int
	}{
		{name: "ok", input: list(1, 2, 3), pred: func(i int) bool { return i < 3 }, want: []int{1, 2}},
		{name: "all elements ignored", input: list(1, 2, 3), pred: func(i int) bool { return i < 0 }, want: []int{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := seq.TakeWhile(tc.input, tc.pred)
			if gv := values(got); !reflect.DeepEqual(gv, tc.want) {
				t.Errorf("result mismatch:\n\twant: %#v\n\t got: %#v", tc.want, gv)
			}
		})
	}
}

func TestZip(t *testing.T) {
	testCases := []struct {
		name string
		as   iter.Seq[int]
		bs   iter.Seq[string]
		want []pair[int, string]
	}{
		{name: "ok", as: list(1, 2, 3), bs: list("a", "b", "c"), want: []pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}},
		{name: "as is shorter than bs", as: list(1, 2), bs: list("a", "b", "c"), want: []pair[int, string]{{1, "a"}, {2, "b"}}},
		{name: "bs is shorter than as", as: list(1, 2, 3), bs: list("a", "b"), want: []pair[int, string]{{1, "a"}, {2, "b"}}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := seq.Zip(tc.as, tc.bs)
			if gv := pairs(got); !reflect.DeepEqual(gv, tc.want) {
				t.Errorf("result mismatch:\n\twant: %#v\n\t got: %#v", tc.want, gv)
			}
		})
	}
}

func TestZipAll(t *testing.T) {
	testCases := []struct {
		name    string
		as      iter.Seq[int]
		bs      iter.Seq[string]
		fillInt int
		fillStr string
		want    []pair[int, string]
	}{
		{name: "ok", as: list(1, 2, 3), bs: list("a", "b", "c"), want: []pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}},
		{name: "as is shorter than bs", fillInt: -1, as: list(1, 2), bs: list("a", "b", "c"), want: []pair[int, string]{{1, "a"}, {2, "b"}, {-1, "c"}}},
		{name: "bs is shorter than as", fillStr: "z", as: list(1, 2, 3), bs: list("a", "b"), want: []pair[int, string]{{1, "a"}, {2, "b"}, {3, "z"}}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := seq.ZipAll(tc.as, tc.bs, tc.fillInt, tc.fillStr)
			if gv := pairs(got); !reflect.DeepEqual(gv, tc.want) {
				t.Errorf("result mismatch:\n\twant: %#v\n\t got: %#v", tc.want, gv)
			}
		})
	}
}

func TestChunk(t *testing.T) {
	testCases := []struct {
		input iter.Seq[string]
		size  int
		want  [][]string
	}{
		{
			input: slices.Values([]string{"a", "b", "c", "d"}),
			size:  2,
			want:  [][]string{{"a", "b"}, {"c", "d"}},
		},
		{
			input: slices.Values([]string{"a", "b", "c", "d", "e"}),
			size:  2,
			want:  [][]string{{"a", "b"}, {"c", "d"}, {"e"}},
		},
		{
			input: slices.Values([]string{"a", "b"}),
			size:  3,
			want:  [][]string{{"a", "b"}},
		},
		{
			input: slices.Values([]string{}),
			size:  2,
			want:  nil,
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("size=%d input=%#v", tc.size, slices.Collect(tc.input)), func(t *testing.T) {
			var got [][]string
			for s := range seq.Chunk(tc.input, tc.size) {
				got = append(got, slices.Collect(s))
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("values:\n\twant: %#v\n\t got: %#v", tc.want, got)
			}
		})
	}
}

func TestUnnestOnce(t *testing.T) {
	testCases := []struct {
		name  string
		input iter.Seq[iter.Seq[int]]
		want  []int
	}{
		{
			name:  "one element",
			input: slices.Values([]iter.Seq[int]{slices.Values([]int{1, 2, 3})}),
			want:  []int{1, 2, 3},
		},
		{
			name:  "multiple elements",
			input: slices.Values([]iter.Seq[int]{slices.Values([]int{1, 2, 3}), slices.Values([]int{4, 5, 6})}),
			want:  []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:  "empty",
			input: slices.Values([]iter.Seq[int]{}),
			want:  nil,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := slices.Collect(seq.UnnestOnce(tc.input))
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("values:\n\twant: %#v\n\t got: %#v", tc.want, got)
			}
		})
	}
}

func TestChunkPairs(t *testing.T) {
	testCases := []struct {
		input iter.Seq[string]
		size  int
		want  [][]string
	}{
		{
			input: slices.Values([]string{"a", "b", "c", "d"}),
			size:  2,
			want:  [][]string{{"a", "b"}, {"c", "d"}},
		},
		{
			input: slices.Values([]string{"a", "b", "c", "d", "e"}),
			size:  2,
			want:  [][]string{{"a", "b"}, {"c", "d"}},
		},
		{
			input: slices.Values([]string{"a", "b"}),
			size:  3,
			want:  [][]string{{"a", "b"}},
		},
		{
			input: slices.Values([]string{}),
			size:  2,
			want:  nil,
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("size=%d input=%#v", tc.size, slices.Collect(tc.input)), func(t *testing.T) {
			var got [][]string
			for a, b := range seq.ChunkPairs(tc.input) {
				got = append(got, []string{a, b})
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("values:\n\twant: %#v\n\t got: %#v", tc.want, got)
			}
		})
	}
}

func TestMap(t *testing.T) {
	testCases := []struct {
		name  string
		input iter.Seq[int]
		fn    func(int) string
		want  []string
	}{
		{
			name:  "int to string",
			input: list(1, 2, 3),
			fn:    func(i int) string { return fmt.Sprintf("num%d", i) },
			want:  []string{"num1", "num2", "num3"},
		},
		{
			name:  "empty input",
			input: list[int](),
			fn:    func(i int) string { return fmt.Sprintf("%d", i) },
			want:  []string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := seq.Map(tc.input, tc.fn)
			if gv := values(got); !reflect.DeepEqual(gv, tc.want) {
				t.Errorf("result mismatch:\n\twant: %#v\n\t got: %#v", tc.want, gv)
			}
		})
	}
}

func TestFlatMap(t *testing.T) {
	testCases := []struct {
		name  string
		input iter.Seq[int]
		fn    func(int) iter.Seq[string]
		want  []string
	}{
		{
			name:  "each element to multiple strings",
			input: list(1, 2, 3),
			fn: func(i int) iter.Seq[string] {
				return list(fmt.Sprintf("a%d", i), fmt.Sprintf("b%d", i))
			},
			want: []string{"a1", "b1", "a2", "b2", "a3", "b3"},
		},
		{
			name:  "some elements produce empty sequences",
			input: list(1, 2, 3),
			fn: func(i int) iter.Seq[string] {
				if i%2 == 0 {
					return list[string]()
				}
				return list(fmt.Sprintf("odd%d", i))
			},
			want: []string{"odd1", "odd3"},
		},
		{
			name:  "empty input",
			input: list[int](),
			fn: func(i int) iter.Seq[string] {
				return list(fmt.Sprintf("%d", i))
			},
			want: []string{},
		},
		{
			name:  "all elements produce empty sequences",
			input: list(1, 2, 3),
			fn: func(i int) iter.Seq[string] {
				return list[string]()
			},
			want: []string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := seq.FlatMap(tc.input, tc.fn)
			if gv := values(got); !reflect.DeepEqual(gv, tc.want) {
				t.Errorf("result mismatch:\n\twant: %#v\n\t got: %#v", tc.want, gv)
			}
		})
	}
}

func TestPairwise(t *testing.T) {
	testCases := []struct {
		name  string
		input []string
		want  []pair[string, string]
	}{
		{
			name:  "empty",
			input: []string{},
			want:  []pair[string, string]{},
		},
		{
			name:  "only one element",
			input: []string{"a"},
			want:  []pair[string, string]{},
		},
		{
			name:  "ok",
			input: []string{"a", "b", "c", "d"},
			want: []pair[string, string]{
				{"a", "b"},
				{"b", "c"},
				{"c", "d"},
			},
		},
		{
			name:  "3 elements",
			input: []string{"a", "b", "c"},
			want: []pair[string, string]{
				{"a", "b"},
				{"b", "c"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := pairs(seq.Pairwise(slices.Values(tc.input)))
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("result mismatch:\n\twant: %#v\n\t got: %#v", tc.want, got)
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

type pair[A, B any] struct {
	A A
	B B
}

func pairs[A, B any](s iter.Seq2[A, B]) []pair[A, B] {
	ret := make([]pair[A, B], 0)
	for a, b := range s {
		ret = append(ret, pair[A, B]{a, b})
	}
	return ret
}
