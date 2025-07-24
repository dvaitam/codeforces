package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func probability(order []byte, R, C int) float64 {
	n := R + C
	L := 2 * n
	dp := make([][]float64, R+1)
	for i := range dp {
		dp[i] = make([]float64, R+1)
	}
	dp[0][0] = 1
	var winA float64
	for step := 0; step < L; step++ {
		next := make([][]float64, R+1)
		for i := range next {
			next[i] = make([]float64, R+1)
		}
		player := order[step]
		remaining := L - step
		for a := 0; a < R; a++ {
			for b := 0; b < R; b++ {
				prob := dp[a][b]
				if prob == 0 {
					continue
				}
				rawUsed := a + b
				rawLeft := 2*R - rawUsed
				pRaw := float64(rawLeft) / float64(remaining)
				pCook := 1 - pRaw
				if player == 'A' {
					if a+1 >= R {
						winA += 0 // A loses
					} else {
						next[a+1][b] += prob * pRaw
					}
					next[a][b] += prob * pCook
				} else {
					if b+1 >= R {
						winA += prob * pRaw
					} else {
						next[a][b+1] += prob * pRaw
					}
					next[a][b] += prob * pCook
				}
			}
		}
		dp = next
	}
	return winA
}

var (
	R, C  int
	S     string
	n     int
	best  float64
	count int64
)

func dfs(pos, a, b int, seq []byte) {
	if pos == len(S) {
		if a == n && b == n {
			pA := probability(seq, R, C)
			diff := math.Abs(pA - (1 - pA))
			if diff < best-1e-12 {
				best = diff
				count = 1
			} else if math.Abs(diff-best) <= 1e-12 {
				count++
			}
		}
		return
	}
	if (S[pos] == 'A' || S[pos] == '?') && a < n {
		seq[pos] = 'A'
		dfs(pos+1, a+1, b, seq)
	}
	if (S[pos] == 'B' || S[pos] == '?') && b < n {
		seq[pos] = 'B'
		dfs(pos+1, a, b+1, seq)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(in, &R, &C); err != nil {
		return
	}
	fmt.Fscan(in, &S)
	n = R + C
	best = 1e9
	seq := make([]byte, len(S))
	dfs(0, 0, 0, seq)
	fmt.Println(count)
}
