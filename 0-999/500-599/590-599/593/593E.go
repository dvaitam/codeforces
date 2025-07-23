package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

type Matrix [][]int64

func makeMatrix(n int) Matrix {
	m := make(Matrix, n)
	for i := range m {
		m[i] = make([]int64, n)
	}
	return m
}

func mulVec(v []int64, m Matrix) []int64 {
	n := len(v)
	res := make([]int64, n)
	for j := 0; j < n; j++ {
		var sum int64
		for i := 0; i < n; i++ {
			if v[i] == 0 || m[i][j] == 0 {
				continue
			}
			sum += v[i] * m[i][j]
			if sum >= MOD*MOD {
				sum %= MOD
			}
		}
		res[j] = sum % MOD
	}
	return res
}

func mulMat(a, b Matrix) Matrix {
	n := len(a)
	res := makeMatrix(n)
	for i := 0; i < n; i++ {
		for k := 0; k < n; k++ {
			if a[i][k] == 0 {
				continue
			}
			aik := a[i][k]
			for j := 0; j < n; j++ {
				if b[k][j] == 0 {
					continue
				}
				res[i][j] = (res[i][j] + aik*b[k][j]) % MOD
			}
		}
	}
	return res
}

func applyPower(v []int64, base Matrix, exp int64) []int64 {
	res := make([]int64, len(v))
	copy(res, v)
	power := base
	for exp > 0 {
		if exp&1 == 1 {
			res = mulVec(res, power)
		}
		exp >>= 1
		if exp > 0 {
			power = mulMat(power, power)
		}
	}
	return res
}

func buildStep(neigh [][]int, maskPrev, maskNext int) Matrix {
	n := len(neigh)
	m := makeMatrix(n)
	for i := 0; i < n; i++ {
		if maskPrev&(1<<i) != 0 {
			continue
		}
		for _, j := range neigh[i] {
			if maskNext&(1<<j) != 0 {
				continue
			}
			m[i][j] = 1
		}
	}
	return m
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m, q int
	if _, err := fmt.Fscan(reader, &n, &m, &q); err != nil {
		return
	}
	total := n * m
	idx := func(x, y int) int { return (x-1)*m + (y - 1) }
	neigh := make([][]int, total)
	for r := 1; r <= n; r++ {
		for c := 1; c <= m; c++ {
			id := idx(r, c)
			neigh[id] = append(neigh[id], id) // stay
			if r > 1 {
				neigh[id] = append(neigh[id], idx(r-1, c))
			}
			if r < n {
				neigh[id] = append(neigh[id], idx(r+1, c))
			}
			if c > 1 {
				neigh[id] = append(neigh[id], idx(r, c-1))
			}
			if c < m {
				neigh[id] = append(neigh[id], idx(r, c+1))
			}
		}
	}

	type event struct {
		tp int
		x  int
		y  int
		t  int64
	}
	events := make([]event, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(reader, &events[i].tp, &events[i].x, &events[i].y, &events[i].t)
	}

	mask := 0
	curTime := int64(1)
	state := make([]int64, total)
	state[idx(1, 1)] = 1
	constMat := buildStep(neigh, mask, mask)

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for _, ev := range events {
		delta := ev.t - curTime
		if ev.tp == 1 { // query
			if delta > 0 {
				state = applyPower(state, constMat, delta)
			}
			pos := idx(ev.x, ev.y)
			fmt.Fprintln(writer, state[pos]%MOD)
			curTime = ev.t
		} else {
			nextMask := mask
			pos := idx(ev.x, ev.y)
			if ev.tp == 2 {
				nextMask |= 1 << pos
			} else {
				nextMask &^= 1 << pos
			}
			if delta > 0 {
				if delta > 1 {
					state = applyPower(state, constMat, delta-1)
				}
				stepMat := buildStep(neigh, mask, nextMask)
				state = mulVec(state, stepMat)
			}
			curTime = ev.t
			mask = nextMask
			constMat = buildStep(neigh, mask, mask)
		}
	}
}
