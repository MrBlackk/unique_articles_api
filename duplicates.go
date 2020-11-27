package main

func isDuplicateArticles(a1 []string, a2 []string) bool {
	d := distance(a1, a2)
	maxLength := max(len(a1), len(a2))
	similarity := 100.0 - float64(d)/float64(maxLength)*100.0
	return similarity >= Conf.ArticleSimilarity
}

// Modified Levenshtein distance for measuring the difference between two articles: a1 and a2.
// Distance between two articles is the minimum number of single-word edits
// (insertions, deletions or substitutions) required to change one article into the other
func distance(a1 []string, a2 []string) int {
	l1 := len(a1)
	l2 := len(a2)
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
			if a1[j-1] != a2[i-1] {
				subst++
			}
			currRow[j] = min(prevRow[j]+1, currRow[j-1]+1, subst)
		}
	}
	return currRow[l1]
}
