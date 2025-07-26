package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, M, K int
	fmt.Fscan(in, &N, &M, &K)
	A := make([][]int, N)
	for i := 0; i < N; i++ {
		A[i] = make([]int, M)
		for j := 0; j < M; j++ {
			fmt.Fscan(in, &A[i][j])
		}
	}

	T := M - K + 1
	if T <= 0 {
		fmt.Fprintln(out, 0)
		return
	}

	bit := make([]int64, T+2)
	update := func(idx int, val int64) {
		for idx <= T {
			if val > bit[idx] {
				bit[idx] = val
			} else {
				// if current value is larger, no need to propagate further
			}
			idx += idx & -idx
		}
	}
	query := func(idx int) int64 {
		var res int64
		for idx > 0 {
			if bit[idx] > res {
				res = bit[idx]
			}
			idx -= idx & -idx
		}
		return res
	}

	for j := 1; j <= M; j++ {
		R := j
		if R > T {
			R = T
		}
		L := j - K + 1
		if L < 1 {
			L = 1
		}
		if L > R {
			continue
		}
		lenj := R - L + 1
		vals := make([]int, N)
		for i := 0; i < N; i++ {
			vals[i] = A[i][j-1]
		}
		sort.Slice(vals, func(a, b int) bool { return vals[a] > vals[b] })
		m := N
		if lenj < m {
			m = lenj
		}
		prefix := make([]int64, m+1)
		for i := 1; i <= m; i++ {
			prefix[i] = prefix[i-1] + int64(vals[i-1])
		}
		for r := 1; r <= m; r++ {
			val := prefix[r]
			for s := L; s <= R-r+1; s++ {
				end := s + r - 1
				if end > T {
					break
				}
				cand := query(s-1) + val
				update(end, cand)
			}
		}
	}

	fmt.Fprintln(out, query(T))
}
