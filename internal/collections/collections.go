package collections

func Chunk[T any](values []T, size int) [][]T {
	if size < 1 {
		return nil
	}
	out := make([][]T, 0, (len(values)+size-1)/size)
	for len(values) > 0 {
		n := size
		if n > len(values) {
			n = len(values)
		}
		part := append([]T(nil), values[:n]...)
		out = append(out, part)
		values = values[n:]
	}
	return out
}
func Map[T any, R any](values []T, fn func(T) R) []R {
	out := make([]R, len(values))
	for i, v := range values {
		out[i] = fn(v)
	}
	return out
}
func Compact[T comparable](values []T, zero T) []T {
	out := make([]T, 0, len(values))
	for _, v := range values {
		if v != zero {
			out = append(out, v)
		}
	}
	return out
}
func Clone[T any](values []T) []T { return append([]T(nil), values...) }
