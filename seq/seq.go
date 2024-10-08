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
