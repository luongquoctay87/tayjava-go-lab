package genericslab

// Unique returns comparable values in first-seen order without duplicates.
func Unique[T comparable](in []T) []T {
	seen := make(map[T]struct{}, len(in))
	out := make([]T, 0, len(in))
	for _, v := range in {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

// First returns the first element or the zero value when the slice is empty.
func First[T any](in []T) (T, bool) {
	var zero T
	if len(in) == 0 {
		return zero, false
	}
	return in[0], true
}

// MapIDs projects a slice through fn (lab Map helper).
func MapIDs[T any, U any](in []T, fn func(T) U) []U {
	out := make([]U, len(in))
	for i, v := range in {
		out[i] = fn(v)
	}
	return out
}
