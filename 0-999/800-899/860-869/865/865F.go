package main

import (
	"fmt"
	"sort"
	"strings"
)

var comb [65][65]int64

func initComb() {
	for i := 0; i <= 60; i++ {
		comb[i][0] = 1
		for j := 1; j <= i; j++ {
			comb[i][j] = comb[i-1][j-1] + comb[i-1][j]
		}
	}
}

func nCr(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return comb[n][k]
}

type State struct {
	val   int64
	count int64
}

func compact(states []State) []State {
	if len(states) <= 1 {
		return states
	}
	sort.Slice(states, func(i, j int) bool {
		return states[i].val < states[j].val
	})
	idx := 0
	for i := 1; i < len(states); i++ {
		if states[i].val == states[idx].val {
			states[idx].count += states[i].count
		} else {
			idx++
			states[idx] = states[i]
		}
	}
	return states[:idx+1]
}

func getDiffAndCount(S string, R, C int, targetDiff int64) (int64, int64) {
	var fwd [32][]State
	fwd[0] = []State{{0, 1}}

	for step := 0; step < R+C; step++ {
		var nextFwd [32][]State
		for a := 0; a <= step; a++ {
			b := step - a
			if a > R+C || b > R+C || len(fwd[a]) == 0 {
				continue
			}

			if S[step] == 'A' || S[step] == '?' {
				if a+1 <= R+C {
					cost := nCr(a, R-1) * nCr(b, R)
					for _, st := range fwd[a] {
						nextFwd[a+1] = append(nextFwd[a+1], State{st.val + cost, st.count})
					}
				}
			}
			if S[step] == 'B' || S[step] == '?' {
				if b+1 <= R+C {
					for _, st := range fwd[a] {
						nextFwd[a] = append(nextFwd[a], State{st.val, st.count})
					}
				}
			}
		}
		for a := 0; a <= step+1; a++ {
			if a <= R+C && (step+1-a) <= R+C {
				nextFwd[a] = compact(nextFwd[a])
			}
		}
		fwd = nextFwd
	}

	var bwd [32][]State
	bwd[R+C] = []State{{0, 1}}

	for step := 2 * (R + C); step > R+C; step-- {
		var nextBwd [32][]State
		for a := 0; a <= R+C; a++ {
			b := step - a
			if b < 0 || b > R+C || len(bwd[a]) == 0 {
				continue
			}

			if S[step-1] == 'A' || S[step-1] == '?' {
				if a-1 >= 0 {
					cost := nCr(a-1, R-1) * nCr(b, R)
					for _, st := range bwd[a] {
						nextBwd[a-1] = append(nextBwd[a-1], State{st.val + cost, st.count})
					}
				}
			}
			if S[step-1] == 'B' || S[step-1] == '?' {
				if b-1 >= 0 {
					for _, st := range bwd[a] {
						nextBwd[a] = append(nextBwd[a], State{st.val, st.count})
					}
				}
			}
		}
		for a := 0; a <= R+C; a++ {
			b := (step - 1) - a
			if b >= 0 && b <= R+C {
				nextBwd[a] = compact(nextBwd[a])
			}
		}
		bwd = nextBwd
	}

	minDiff := int64(-1)
	K := nCr(R+C, R) * nCr(R+C, R)

	for a := 0; a <= R+C; a++ {
		fw := fwd[a]
		bw := bwd[a]
		if len(fw) == 0 || len(bw) == 0 {
			continue
		}
		for i := 0; i < len(fw); i++ {
			target := (K + 1) / 2 - fw[i].val
			idx := sort.Search(len(bw), func(j int) bool {
				return bw[j].val >= target
			})
			if idx < len(bw) {
				d := fw[i].val + bw[idx].val
				d = d*2 - K
				if d < 0 {
					d = -d
				}
				if minDiff == -1 || d < minDiff {
					minDiff = d
				}
			}
			if idx > 0 {
				d := fw[i].val + bw[idx-1].val
				d = d*2 - K
				if d < 0 {
					d = -d
				}
				if minDiff == -1 || d < minDiff {
					minDiff = d
				}
			}
		}
	}

	if targetDiff != -1 {
		minDiff = targetDiff
	}

	totalCount := int64(0)
	if minDiff != -1 {
		for a := 0; a <= R+C; a++ {
			fw := fwd[a]
			bw := bwd[a]
			if len(fw) == 0 || len(bw) == 0 {
				continue
			}
			for i := 0; i < len(fw); i++ {
				var targets []int64
				if (K+minDiff)%2 == 0 {
					targets = append(targets, (K+minDiff)/2-fw[i].val)
				}
				if minDiff != 0 && (K-minDiff)%2 == 0 {
					targets = append(targets, (K-minDiff)/2-fw[i].val)
				}

				for _, t := range targets {
					idx := sort.Search(len(bw), func(j int) bool {
						return bw[j].val >= t
					})
					if idx < len(bw) && bw[idx].val == t {
						totalCount += fw[i].count * bw[idx].count
					}
				}
			}
		}
	}

	return minDiff, totalCount
}

func main() {
	var R, C int
	if _, err := fmt.Scan(&R, &C); err != nil {
		return
	}
	var S string
	if _, err := fmt.Scan(&S); err != nil {
		return
	}

	initComb()

	qMarks := strings.Repeat("?", 2*(R+C))
	globalMinDiff, _ := getDiffAndCount(qMarks, R, C, -1)
	_, ans := getDiffAndCount(S, R, C, globalMinDiff)

	fmt.Println(ans)
}
