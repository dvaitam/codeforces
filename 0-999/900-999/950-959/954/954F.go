package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 1000000007

type Matrix [3][3]int64

func matMul(a, b Matrix) Matrix {
	var c Matrix
	for i := 0; i < 3; i++ {
		for k := 0; k < 3; k++ {
			if a[i][k] == 0 {
				continue
			}
			for j := 0; j < 3; j++ {
				c[i][j] = (c[i][j] + a[i][k]*b[k][j]) % mod
			}
		}
	}
	return c
}

func matPow(m Matrix, p int64) Matrix {
	var res Matrix
	for i := 0; i < 3; i++ {
		res[i][i] = 1
	}
	for p > 0 {
		if p&1 == 1 {
			res = matMul(m, res)
		}
		m = matMul(m, m)
		p >>= 1
	}
	return res
}

func matVecMul(m Matrix, v [3]int64) [3]int64 {
	var r [3]int64
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			r[i] = (r[i] + m[i][j]*v[j]) % mod
		}
	}
	return r
}

func baseMatrix(mask int) Matrix {
	var m Matrix
	for from := 0; from < 3; from++ {
		if mask&(1<<from) != 0 {
			continue
		}
		for d := -1; d <= 1; d++ {
			to := from + d
			if to < 0 || to >= 3 {
				continue
			}
			if mask&(1<<to) != 0 {
				continue
			}
			m[to][from] = (m[to][from] + 1) % mod
		}
	}
	return m
}

func stepMatrix(prev, next int) Matrix {
	var m Matrix
	for from := 0; from < 3; from++ {
		if prev&(1<<from) != 0 {
			continue
		}
		for d := -1; d <= 1; d++ {
			to := from + d
			if to < 0 || to >= 3 {
				continue
			}
			if next&(1<<to) != 0 {
				continue
			}
			m[to][from] = (m[to][from] + 1) % mod
		}
	}
	return m
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var mval int64
	if _, err := fmt.Fscan(in, &n, &mval); err != nil {
		return
	}

	events := make(map[int64][][2]int)
	for i := 0; i < n; i++ {
		var a int
		var l, r int64
		fmt.Fscan(in, &a, &l, &r)
		events[l] = append(events[l], [2]int{a - 1, 1})
		events[r+1] = append(events[r+1], [2]int{a - 1, -1})
	}

	var positions []int64
	for x := range events {
		if x <= mval {
			positions = append(positions, x)
		}
	}
	sort.Slice(positions, func(i, j int) bool { return positions[i] < positions[j] })

	var base [8]Matrix
	for mask := 0; mask < 8; mask++ {
		base[mask] = baseMatrix(mask)
	}

	var dp [3]int64
	dp[1] = 1
	var cnt [3]int
	mask := 0
	current := int64(1)
	idx := 0
	for idx < len(positions) {
		x := positions[idx]
		if x-current-1 > 0 {
			dp = matVecMul(matPow(base[mask], x-current-1), dp)
			current = x - 1
		}
		prevMask := mask
		for idx < len(positions) && positions[idx] == x {
			for _, ev := range events[x] {
				cnt[ev[0]] += ev[1]
				if cnt[ev[0]] > 0 {
					mask |= 1 << ev[0]
				} else {
					mask &^= 1 << ev[0]
				}
			}
			idx++
		}
		dp = matVecMul(stepMatrix(prevMask, mask), dp)
		current = x
	}
	if current < mval {
		dp = matVecMul(matPow(base[mask], mval-current), dp)
	}

	fmt.Fprintln(out, dp[1]%mod)
}
