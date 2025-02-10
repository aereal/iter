package seq2

import "iter"

// Flip swaps the elements of a given sequence of pairs.
func Flip[T any, U any](input iter.Seq2[T, U]) iter.Seq2[U, T] {
	return func(yield func(U, T) bool) {
		for t, u := range input {
			if !yield(u, t) {
				return
			}
		}
	}
}
