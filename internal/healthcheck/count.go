package healthcheck

func Count(results []Result) (int, int) {
	ok := 0
	failed := 0
	for _, result := range results {
		if result.OK {
			ok++
		} else {
			failed++
		}
	}
	return ok, failed
}

func Names(results []Result) []string {
	names := make([]string, 0, len(results))
	for _, result := range results {
		names = append(names, result.Name)
	}
	return names
}
