package seq

import "iter"

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
