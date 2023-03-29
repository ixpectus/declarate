package tools

func Filter[T any](slice []T, f func(T) bool) []T {
	var res []T
	for _, e := range slice {
		if f(e) {
			res = append(res, e)
		}
	}

	return res
}

func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}

	return false
}

func Intersect[T comparable](a, b []T) []T {
	set := make([]T, 0)
	for _, v := range a {
		if Contains(b, v) {
			set = append(set, v)
		}
	}

	return set
}
