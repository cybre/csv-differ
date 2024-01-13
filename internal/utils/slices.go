package utils

func Map[S ~[]E, E, R any](s S, f func(E) R) []R {
	result := make([]R, len(s))
	for i, v := range s {
		result[i] = f(v)
	}

	return result
}
