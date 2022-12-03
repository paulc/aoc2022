package util

func Must[T any](r T, err error) T {
	if err != nil {
		panic(err)
	}
	return r
}

func Take[T any](s []T, n int) (out [][]T) {
	i := 0
	group := []T{}
	for _, v := range s {
		if i != 0 && i%n == 0 {
			out = append(out, group)
			group = []T{}
		}
		group = append(group, v)
		i++
	}
	if len(group) > 0 {
		out = append(out, group)
	}
	return
}
