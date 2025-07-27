package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	a := make([]int, n)
	b := make([]int, m)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for j := 0; j < m; j++ {
		fmt.Fscan(in, &b[j])
	}
	c := make([][]int, n)
	for i := 0; i < n; i++ {
		c[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &c[i][j])
		}
	}

	sumA, sumB := 0, 0
	for _, v := range a {
		sumA += v
	}
	for _, v := range b {
		sumB += v
	}
	if sumA > sumB {
		fmt.Fprintln(out, -1)
		return
	}

	pow := make([]int, n)
	pow[0] = 1
	for i := 1; i < n; i++ {
		pow[i] = pow[i-1] * 5
	}
	finalNeed := 0
	for i := 0; i < n; i++ {
		finalNeed += a[i] * pow[i]
	}

	type state struct {
		need int
		v1   int
		v2   int
		rem  int
	}

	start := state{0, 0, 0, 0}
	dp := map[state]int{start: 0}
	best := int64(1<<63 - 1)

	for len(dp) > 0 {
		next := make(map[state]int)
		for st, cost := range dp {
			if int64(cost) >= best {
				continue
			}
			if st.v2 == m {
				if st.need == finalNeed && int64(cost) < best {
					best = int64(cost)
				}
				continue
			}
			i := st.v1
			j := st.v2
			cur := (st.need / pow[i]) % 5
			maxF := a[i] - cur
			if maxF > b[j]-st.rem {
				maxF = b[j] - st.rem
			}
			for f := 0; f <= maxF; f++ {
				nn := st.need + f*pow[i]
				nc := cost
				if f > 0 {
					nc += c[i][j]
				}
				nv1 := i + 1
				nv2 := j
				nr := st.rem + f
				if nv1 == n {
					nv1 = 0
					nv2 = j + 1
					nr = 0
				}
				nst := state{nn, nv1, nv2, nr}
				if nv2 == m {
					if nn == finalNeed && int64(nc) < best {
						best = int64(nc)
					}
					continue
				}
				if prev, ok := next[nst]; !ok || nc < prev {
					next[nst] = nc
				}
			}
		}
		dp = next
	}

	if best == int64(1<<63-1) {
		fmt.Fprintln(out, -1)
	} else {
		fmt.Fprintln(out, best)
	}
}
