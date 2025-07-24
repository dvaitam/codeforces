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
	var n int
	var s string
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	fmt.Fscan(in, &s)

	pos := make([][]int, 3)
	prefix := make([][]int, 3)
	for i := 0; i < 3; i++ {
		prefix[i] = make([]int, n+1)
	}

	for i := 0; i < n; i++ {
		for t := 0; t < 3; t++ {
			prefix[t][i+1] = prefix[t][i]
		}
		var cat int
		if s[i] == 'V' {
			cat = 0
		} else if s[i] == 'K' {
			cat = 1
		} else {
			cat = 2
		}
		pos[cat] = append(pos[cat], i)
		prefix[cat][i+1]++
	}

	cntV, cntK, cntO := len(pos[0]), len(pos[1]), len(pos[2])
	size0, size1, size2 := cntV+1, cntK+1, cntO+1

	const INF = int(1e9)
	dp := make([]int, size0*size1*size2*4)
	for i := range dp {
		dp[i] = INF
	}
	idx := func(a, b, c, last int) int {
		return ((a*size1+b)*size2+c)*4 + last
	}
	dp[idx(0, 0, 0, 3)] = 0

	for a := 0; a <= cntV; a++ {
		for b := 0; b <= cntK; b++ {
			for c := 0; c <= cntO; c++ {
				for last := 0; last < 4; last++ {
					cur := dp[idx(a, b, c, last)]
					if cur == INF {
						continue
					}
					for t := 0; t < 3; t++ {
						if last == 0 && t == 1 {
							continue // forbid "VK"
						}
						var posIdx int
						if t == 0 {
							if a >= cntV {
								continue
							}
							posIdx = pos[0][a]
						} else if t == 1 {
							if b >= cntK {
								continue
							}
							posIdx = pos[1][b]
						} else {
							if c >= cntO {
								continue
							}
							posIdx = pos[2][c]
						}
						used := min(a, prefix[0][posIdx]) +
							min(b, prefix[1][posIdx]) +
							min(c, prefix[2][posIdx])
						cost := posIdx - used
						na, nb, nc := a, b, c
						if t == 0 {
							na++
						} else if t == 1 {
							nb++
						} else {
							nc++
						}
						id2 := idx(na, nb, nc, t)
						if cur+cost < dp[id2] {
							dp[id2] = cur + cost
						}
					}
				}
			}
		}
	}

	ans := INF
	for last := 0; last < 3; last++ {
		if dp[idx(cntV, cntK, cntO, last)] < ans {
			ans = dp[idx(cntV, cntK, cntO, last)]
		}
	}
	fmt.Println(ans)
}
