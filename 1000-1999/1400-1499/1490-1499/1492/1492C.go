package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	var s, t string
	fmt.Fscan(reader, &s)
	fmt.Fscan(reader, &t)

	pref := make([]int, m)
	suf := make([]int, m)

	// compute earliest positions for each prefix of t
	j := 0
	for i := 0; i < n && j < m; i++ {
		if s[i] == t[j] {
			pref[j] = i
			j++
		}
	}

	// compute latest positions for each suffix of t
	j = m - 1
	for i := n - 1; i >= 0 && j >= 0; i-- {
		if s[i] == t[j] {
			suf[j] = i
			j--
		}
	}

	ans := 0
	for i := 0; i < m-1; i++ {
		diff := suf[i+1] - pref[i]
		if diff > ans {
			ans = diff
		}
	}
	fmt.Println(ans)
}
