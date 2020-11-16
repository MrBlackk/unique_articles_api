package main

func min(vars ...int) int {
	min := vars[0]
	for _, v := range vars {
		if min > v {
			min = v
		}
	}
	return min
}

func max(vars ...int) int {
	max := vars[0]
	for _, v := range vars {
		if max < v {
			max = v
		}
	}
	return max
}
