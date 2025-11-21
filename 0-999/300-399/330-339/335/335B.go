package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(in, &s)
	n := len(s)

	const alphabet = 26
	nextPos := make([][alphabet]int, n+1)
	for c := 0; c < alphabet; c++ {
		nextPos[n][c] = n
	}
	for i := n - 1; i >= 0; i-- {
		nextPos[i] = nextPos[i+1]
		nextPos[i][s[i]-'a'] = i
	}

	prevPos := make([][alphabet]int, n+1)
	for c := 0; c < alphabet; c++ {
		prevPos[0][c] = -1
	}
	for i := 0; i < n; i++ {
		prevPos[i+1] = prevPos[i]
		prevPos[i+1][s[i]-'a'] = i
	}

	maxPairs := 50
	ans := make([]byte, 0, maxPairs)
	l, r := 0, n-1
	for len(ans) < maxPairs && l < r {
		found := false
		for c := 0; c < alphabet; c++ {
			left := nextPos[l][c]
			if left >= r {
				continue
			}
			right := prevPos[r+1][c]
			if left < right {
				ans = append(ans, byte('a'+c))
				l = left + 1
				r = right - 1
				found = true
				break
			}
		}
		if !found {
			break
		}
	}

	res := make([]byte, 0, len(ans)*2+1)
	res = append(res, ans...)
	if len(res) < 100 && l <= r {
		res = append(res, s[l])
	}
	for i := len(ans) - 1; i >= 0; i-- {
		res = append(res, ans[i])
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintln(out, string(res))
}
