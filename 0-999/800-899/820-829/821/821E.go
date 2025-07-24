package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 1000000007

type event struct {
	c   int
	add bool
}

type matrix [][]int64

func matMul(a, b matrix) matrix {
	n := len(a)
	res := make(matrix, n)
	for i := range res {
		res[i] = make([]int64, n)
	}
	for i := 0; i < n; i++ {
		for k := 0; k < n; k++ {
			if a[i][k] == 0 {
				continue
			}
			for j := 0; j < n; j++ {
				if b[k][j] == 0 {
					continue
				}
				res[i][j] = (res[i][j] + a[i][k]*b[k][j]) % mod
			}
		}
	}
	return res
}

func matPow(m matrix, p int64) matrix {
	n := len(m)
	res := make(matrix, n)
	for i := range res {
		res[i] = make([]int64, n)
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

func matVecMul(m matrix, v []int64) []int64 {
	n := len(m)
	r := make([]int64, n)
	for i := 0; i < n; i++ {
		var sum int64
		for j := 0; j < n; j++ {
			if m[i][j] == 0 || v[j] == 0 {
				continue
			}
			sum = (sum + m[i][j]*v[j]) % mod
		}
		r[i] = sum
	}
	return r
}

func baseMatrix(maxY int) matrix {
	n := maxY + 1
	m := make(matrix, n)
	for i := range m {
		m[i] = make([]int64, n)
	}
	for from := 0; from < n; from++ {
		for d := -1; d <= 1; d++ {
			to := from + d
			if to >= 0 && to < n {
				m[to][from] = (m[to][from] + 1) % mod
			}
		}
	}
	return m
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	events := make(map[int64][]event)
	positionsSet := map[int64]struct{}{0: {}, k: {}}

	for i := 0; i < n; i++ {
		var a, b int64
		var c int
		fmt.Fscan(reader, &a, &b, &c)
		events[a] = append(events[a], event{c, true})
		events[b+1] = append(events[b+1], event{c, false})
		positionsSet[a] = struct{}{}
		positionsSet[b+1] = struct{}{}
	}

	positions := make([]int64, 0, len(positionsSet))
	for p := range positionsSet {
		positions = append(positions, p)
	}
	sort.Slice(positions, func(i, j int) bool { return positions[i] < positions[j] })

	base := make([]matrix, 16)
	for i := 0; i <= 15; i++ {
		base[i] = baseMatrix(i)
	}

	dp := make([]int64, 16)
	dp[0] = 1
	cnt := make([]int, 16)

	for idx := 0; idx < len(positions)-1; idx++ {
		pos := positions[idx]
		nextPos := positions[idx+1]

		for _, ev := range events[pos] {
			if ev.add {
				cnt[ev.c]++
			} else {
				cnt[ev.c]--
			}
		}

		if pos >= k {
			break
		}

		limit := -1
		for i := 0; i <= 15; i++ {
			if cnt[i] > 0 {
				limit = i
				break
			}
		}
		if limit == -1 {
			limit = 15
		}

		for i := limit + 1; i < 16; i++ {
			dp[i] = 0
		}

		end := nextPos
		if end > k {
			end = k
		}
		stepLen := end - pos
		if stepLen > 0 {
			mat := matPow(base[limit], stepLen)
			tmp := matVecMul(mat, dp[:limit+1])
			for i := 0; i <= limit; i++ {
				dp[i] = tmp[i]
			}
		}
	}

	if evs, ok := events[k]; ok {
		for _, ev := range evs {
			if ev.add {
				cnt[ev.c]++
			} else {
				cnt[ev.c]--
			}
		}
	}
	limit := -1
	for i := 0; i <= 15; i++ {
		if cnt[i] > 0 {
			limit = i
			break
		}
	}
	if limit == -1 {
		limit = 15
	}
	for i := limit + 1; i < 16; i++ {
		dp[i] = 0
	}

	fmt.Fprintln(writer, dp[0]%mod)
}
