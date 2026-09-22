package algo

// Levenshtein returns the edit distance between a and b.
func Levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	n, m := len(ra), len(rb)
	if n == 0 {
		return m
	}
	if m == 0 {
		return n
	}
	prev := make([]int, m+1)
	curr := make([]int, m+1)
	for j := 0; j <= m; j++ {
		prev[j] = j
	}
	for i := 1; i <= n; i++ {
		curr[0] = i
		for j := 1; j <= m; j++ {
			cost := 0
			if ra[i-1] != rb[j-1] {
				cost = 1
			}
			del := prev[j] + 1
			ins := curr[j-1] + 1
			sub := prev[j-1] + cost
			curr[j] = min3(del, ins, sub)
		}
		prev, curr = curr, prev
	}
	return prev[m]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// LCSLength returns the length of the longest common subsequence of a and b.
func LCSLength(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	n, m := len(ra), len(rb)
	if n == 0 || m == 0 {
		return 0
	}
	prev := make([]int, m+1)
	curr := make([]int, m+1)
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if ra[i-1] == rb[j-1] {
				curr[j] = prev[j-1] + 1
			} else if prev[j] >= curr[j-1] {
				curr[j] = prev[j]
			} else {
				curr[j] = curr[j-1]
			}
		}
		prev, curr = curr, prev
		for j := range curr {
			curr[j] = 0
		}
	}
	return prev[m]
}

// KMPIndex returns the first index of pattern in text, or -1.
func KMPIndex(text, pattern string) int {
	if pattern == "" {
		return 0
	}
	t, p := []rune(text), []rune(pattern)
	if len(p) > len(t) {
		return -1
	}
	lps := buildLPS(p)
	i, j := 0, 0
	for i < len(t) {
		if t[i] == p[j] {
			i++
			j++
			if j == len(p) {
				return i - j
			}
			continue
		}
		if j > 0 {
			j = lps[j-1]
		} else {
			i++
		}
	}
	return -1
}

func buildLPS(p []rune) []int {
	lps := make([]int, len(p))
	length := 0
	i := 1
	for i < len(p) {
		if p[i] == p[length] {
			length++
			lps[i] = length
			i++
			continue
		}
		if length > 0 {
			length = lps[length-1]
			continue
		}
		lps[i] = 0
		i++
	}
	return lps
}
