package seq

import (
	"iter"
	"slices"
)

// Drop returns new [iter.Seq] that yields all elements from the argument except first n ones.
//
// If n < 0, [iter.Seq] that yields all elements from the argument.
func Drop[T any](s iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		var seen int
		for v := range s {
			seen++
			if seen <= n {
				continue
			}
			if !yield(v) {
				break
			}
		}
	}
}

// DropWhile returns new [iter.Seq] that yields the elements but dropped longest prefix satisfy the predicate from the argument.
func DropWhile[T any](s iter.Seq[T], predicate func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		skip := true
		for v := range s {
			if skip && predicate(v) {
				continue
			}
			skip = false
			if !yield(v) {
				break
			}
		}
	}
}

// Take returns new [iter.Seq] that yields first n elements from the argument.
func Take[T any](s iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		var seen int
		for v := range s {
			seen++
			if seen > n {
				break
			}
			if !yield(v) {
				break
			}
		}
	}
}

// TakeWhile return new [iter.Seq] that yields longest prefix of elements from the argument.
func TakeWhile[T any](s iter.Seq[T], predicate func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if !predicate(v) {
				break
			}
			if !yield(v) {
				break
			}
		}
	}
}

// Zip returns new [iter.Seq2] yields elements that by combining corresponding elements in pairs from each sequences.
func Zip[A, B any](as iter.Seq[A], bs iter.Seq[B]) iter.Seq2[A, B] {
	return func(yield func(A, B) bool) {
		nextA, stopA := iter.Pull(as)
		defer stopA()
		nextB, stopB := iter.Pull(bs)
		defer stopB()
		for {
			a, ok := nextA()
			if !ok {
				break
			}
			b, ok := nextB()
			if !ok {
				break
			}
			if !yield(a, b) {
				break
			}
		}
	}
}

func ZipAll[A, B any](as iter.Seq[A], bs iter.Seq[B], fillA A, fillB B) iter.Seq2[A, B] {
	return func(yield func(A, B) bool) {
		nextA, stopA := iter.Pull(as)
		defer stopA()
		nextB, stopB := iter.Pull(bs)
		defer stopB()
		for {
			a, foundA := nextA()
			if !foundA {
				a = fillA
			}
			b, foundB := nextB()
			if !foundB {
				b = fillB
			}
			if !foundA && !foundB {
				break
			}
			if !yield(a, b) {
				break
			}
		}

	}
}

// Chunk returns an iterator over consecutive elements of up to n elements of s.
func Chunk[T any](s iter.Seq[T], n int) iter.Seq[iter.Seq[T]] {
	return func(yield func(iter.Seq[T]) bool) {
		buf := make([]T, 0, n)
		for el := range s {
			buf = append(buf, el)
			if len(buf) >= n {
				if !yield(slices.Values(buf)) {
					break
				}
				buf = buf[0:0]
			}
		}
		if len(buf) > 0 {
			_ = yield(slices.Values(buf))
		}
	}
}

// UnnestOnce flattens a sequence of sequences into a single sequence.
func UnnestOnce[T any](nested iter.Seq[iter.Seq[T]]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for xs := range nested {
			for x := range xs {
				if !yield(x) {
					return
				}
			}
		}
	}
}

// ChunkPairs returns an iterator iter.Seq2[T, T] that generates consecutive pairs (T, T) from the input iterator iter.Seq[T].
// If the input sequence has an odd number of elements, the last element is ignored.
func ChunkPairs[T any](input iter.Seq[T]) iter.Seq2[T, T] {
	return func(yield func(T, T) bool) {
		it, stop := iter.Pull(input)
		defer stop()
		for {
			valFirst, okFirst := it()
			if !okFirst {
				return
			}
			valSecond, okSecond := it()
			if !okSecond {
				return
			}
			if !yield(valFirst, valSecond) {
				return
			}
		}
	}
}

// Map transforms each element by applying the given function.
func Map[T any, R any](s iter.Seq[T], f func(T) R) iter.Seq[R] {
	return func(yield func(R) bool) {
		for t := range s {
			if !yield(f(t)) {
				return
			}
		}
	}
}

// FlatMap transforms each element into a sequence and flattens the result.
func FlatMap[T any, R any](s iter.Seq[T], f func(T) iter.Seq[R]) iter.Seq[R] {
	return func(yield func(R) bool) {
		for t := range s {
			for r := range f(t) {
				if !yield(r) {
					return
				}
			}
		}
	}
}
