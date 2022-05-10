package util

func ValueSlice[K comparable, V any](m map[K]V) []V {
	ret := make([]V, 0, len(m))
	for _, val := range m {
		ret = append(ret, val)
	}
	return ret
}
