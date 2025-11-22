package main

import (
	"bufio"
	"fmt"
	"os"
)

// Precompute LCP for all suffix pairs in O(n^2).
func buildLCP(s string) [][]int {
	n := len(s)
	lcp := make([][]int, n+1)
	for i := range lcp {
		lcp[i] = make([]int, n+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			if s[i] == s[j] {
				lcp[i][j] = 1 + lcp[i+1][j+1]
			}
		}
	}
	return lcp
}

// lexLess tests whether substring s[a:b] is lexicographically smaller than s[c:d].
func lexLess(s string, lcp [][]int, a, b, c, d int) bool {
	len1 := b - a
	len2 := d - c
	common := lcp[a][c]
	if common >= len1 || common >= len2 {
		return len1 < len2
	}
	return s[a+common] < s[c+common]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	n := len(s)
	lcp := buildLCP(s)

	bestCnt := make([]int, n+1)   // maximum words for prefix [0:i)
	lastStart := make([]int, n+1) // start index of the last word for that prefix
	prevPos := make([]int, n+1)   // previous boundary to reconstruct
	for i := 1; i <= n; i++ {
		bestCnt[i] = -1
		lastStart[i] = -1
		prevPos[i] = -1
	}
	bestCnt[0] = 0
	lastStart[0] = -1

	for i := 0; i < n; i++ {
		if bestCnt[i] == -1 {
			continue
		}
		for j := i + 1; j <= n; j++ {
			if i == 0 || lexLess(s, lcp, lastStart[i], i, i, j) {
				newCnt := bestCnt[i] + 1
				if newCnt > bestCnt[j] {
					bestCnt[j] = newCnt
					lastStart[j] = i
					prevPos[j] = i
				} else if newCnt == bestCnt[j] {
					if lexLess(s, lcp, i, j, lastStart[j], j) {
						lastStart[j] = i
						prevPos[j] = i
					}
				}
			}
		}
	}

	k := bestCnt[n]
	fmt.Fprintln(out, k)
	words := make([]string, 0, k)
	pos := n
	for pos > 0 {
		st := lastStart[pos]
		words = append(words, s[st:pos])
		pos = prevPos[pos]
	}
	for i := k - 1; i >= 0; i-- {
		fmt.Fprintln(out, words[i])
	}
}
