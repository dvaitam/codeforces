package main

import (
	"bufio"
	"fmt"
	"os"
)

const negInf = int(-1e9)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)

	if n == 0 {
		fmt.Println(0)
		return
	}
	if n == 1 {
		fmt.Println(1)
		return
	}

	d := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		if s[i] != s[i+1] {
			d[i] = 1
		}
	}

	L := n - k + 1
	bestPerResidue := make([][]int, k)
	for r := 1; r <= k; r++ {
		bestPerResidue[r-1] = solveResidue(r, n, m, k, L, d)
	}

	global := make([]int, m+1)
	for i := range global {
		global[i] = negInf
	}
	global[0] = 0

	for _, best := range bestPerResidue {
		next := make([]int, m+1)
		for i := range next {
			next[i] = negInf
		}
		for used := 0; used <= m; used++ {
			if global[used] == negInf {
				continue
			}
			for c := 0; c < len(best); c++ {
				if best[c] == negInf {
					continue
				}
				if used+c > m {
					continue
				}
				val := global[used] + best[c]
				if val > next[used+c] {
					next[used+c] = val
				}
			}
		}
		global = next
	}

	ans := 0
	for used := 0; used <= m; used++ {
		if global[used] > ans {
			ans = global[used]
		}
	}
	fmt.Println(ans + 1)
}

func solveResidue(r, n, m, k, L int, d []int) []int {
	if r > n {
		return []int{0}
	}
	qMax := (n - r) / k
	nodeCnt := qMax + 2

	varNodes := 0
	for j := r; j <= L; j += k {
		varNodes++
	}
	maxOps := varNodes
	if maxOps > m {
		maxOps = m
	}

	dpPrev := make([][2]int, maxOps+1)
	for i := range dpPrev {
		dpPrev[i][0] = negInf
		dpPrev[i][1] = negInf
	}
	dpPrev[0][0] = 0

	qStart := 0
	if r == 1 {
		qStart = 1
	}

	for p := 1; p <= nodeCnt-1; p++ {
		dpCur := make([][2]int, maxOps+1)
		for i := range dpCur {
			dpCur[i][0] = negInf
			dpCur[i][1] = negInf
		}
		j := r + (p-1)*k
		isVar := j >= 1 && j <= L
		allowed := []int{0}
		if isVar {
			allowed = append(allowed, 1)
		}
		for used := 0; used <= maxOps; used++ {
			for prevVal := 0; prevVal <= 1; prevVal++ {
				prevScore := dpPrev[used][prevVal]
				if prevScore == negInf {
					continue
				}
				for _, val := range allowed {
					cost := 0
					if isVar && val == 1 {
						cost = 1
					}
					newUsed := used + cost
					if newUsed > maxOps {
						continue
					}
					score := prevScore
					q := p - 1
					if q >= qStart {
						t := r + q*k
						edgeIdx := t - 2
						if edgeIdx >= 0 {
							diff := prevVal ^ val
							if diff != d[edgeIdx] {
								score++
							}
						}
					}
					if score > dpCur[newUsed][val] {
						dpCur[newUsed][val] = score
					}
				}
			}
		}
		dpPrev = dpCur
	}

	best := make([]int, maxOps+1)
	for used := 0; used <= maxOps; used++ {
		best[used] = maxInt(dpPrev[used][0], dpPrev[used][1])
	}
	return best
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
