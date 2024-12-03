package sutando

func purge(a []any) any {
	if len(a) == 1 {
		return a[0]
	}
	return a
}
