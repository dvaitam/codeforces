package main

import (
	"bufio"
	"fmt"
	"os"
)

func buildPi(t []byte) []int {
	m := len(t)
	pi := make([]int, m)
	for i := 1; i < m; i++ {
		j := pi[i-1]
		for j > 0 && t[i] != t[j] {
			j = pi[j-1]
		}
		if t[i] == t[j] {
			j++
		}
		pi[i] = j
	}
	return pi
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s, t string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	fmt.Fscan(reader, &t)

	sb := []byte(s)
	tb := []byte(t)
	m := len(tb)
	pi := buildPi(tb)

	const negInf int = -1000000000
	dp := make([]int, m)
	for i := range dp {
		dp[i] = negInf
	}
	dp[0] = 0

	// precompute transitions
	nxt := make([][26]int, m)
	add := make([][26]int, m)
	for j := 0; j < m; j++ {
		for c := 0; c < 26; c++ {
			x := j
			for x > 0 && tb[x] != byte('a'+c) {
				x = pi[x-1]
			}
			if tb[x] == byte('a'+c) {
				x++
			}
			if x == m {
				add[j][c] = 1
				x = pi[m-1]
			}
			nxt[j][c] = x
		}
	}

	for _, ch := range sb {
		newdp := make([]int, m)
		for i := range newdp {
			newdp[i] = negInf
		}
		if ch == '?' {
			for j := 0; j < m; j++ {
				if dp[j] == negInf {
					continue
				}
				val := dp[j]
				for c := 0; c < 26; c++ {
					ns := nxt[j][c]
					nv := val + add[j][c]
					if nv > newdp[ns] {
						newdp[ns] = nv
					}
				}
			}
		} else {
			idx := int(ch - 'a')
			for j := 0; j < m; j++ {
				if dp[j] == negInf {
					continue
				}
				ns := nxt[j][idx]
				nv := dp[j] + add[j][idx]
				if nv > newdp[ns] {
					newdp[ns] = nv
				}
			}
		}
		dp = newdp
	}

	ans := 0
	for _, v := range dp {
		if v > ans {
			ans = v
		}
	}

	fmt.Fprintln(writer, ans)
}
