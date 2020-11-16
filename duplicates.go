package main

func isDuplicate(w1 []string, w2 []string) bool {
	d := distance(w1, w2)
	p := 100.0 - float64(d)/float64(max(len(w1), len(w2)))*100.0
	return p >= float64(Conf.ArticleSimilarity)
}

func distance(s1 []string, s2 []string) int {
	l1 := len(s1)
	l2 := len(s2)
	currRow := make([]int, l1+1)
	prevRow := make([]int, l1+1)
	for i := range currRow {
		currRow[i] = i
	}
	for i := 1; i <= l2; i++ {
		for j := range currRow {
			prevRow[j] = currRow[j]
			if j == 0 {
				currRow[j] = i
				continue
			}
			subst := prevRow[j-1]
			if s1[j-1] != s2[i-1] {
				subst++
			}
			currRow[j] = min(prevRow[j]+1, currRow[j-1]+1, subst)
		}
	}
	return currRow[l1]
}
