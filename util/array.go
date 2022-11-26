package util

func ArrayTranspose[T any](in [][]T) [][]T {
	if len(in) == 0 {
		return in
	}
	w := len(in[0])
	h := len(in)
	out := make([][]T, w)
	for x := 0; x < w; x++ {
		out[x] = make([]T, h)
	}
	for x := 0; x < h; x++ {
		for y := 0; y < w; y++ {
			out[y][x] = in[x][y]
		}
	}
	return out
}
