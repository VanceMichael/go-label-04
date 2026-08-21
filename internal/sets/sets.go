package sets

func Difference[T comparable](a, b []T) []T {
	seen := map[T]bool{}
	for _, v := range b {
		seen[v] = true
	}
	out := make([]T, 0)
	for _, v := range a {
		if !seen[v] {
			out = append(out, v)
		}
	}
	return out
}
func Union[T comparable](a, b []T) []T {
	seen := map[T]bool{}
	out := make([]T, 0)
	for _, v := range append(append([]T{}, a...), b...) {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}
func Contains[T comparable](values []T, target T) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}
