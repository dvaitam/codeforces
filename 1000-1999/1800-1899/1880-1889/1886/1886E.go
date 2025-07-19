package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n, m int
	fmt.Fscan(reader, &n, &m)

	a := make([]struct{ val, idx int }, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i].val)
		a[i].idx = i
	}
	sort.Slice(a, func(i, j int) bool { return a[i].val < a[j].val })
	b := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &b[i])
	}
	// prepare d arrays
	d := make([][]int, m)
	for i := 0; i < m; i++ {
		d[i] = make([]int, n)
		for j := 0; j < n; j++ {
			d[i][j] = -1
		}
		for j := 0; j < n; j++ {
			tmp := j + (b[i]+a[j].val-1)/a[j].val - 1
			if tmp >= n {
				continue
			}
			if d[i][tmp] < j {
				d[i][tmp] = j
			}
		}
		for j := 0; j+1 < n; j++ {
			if d[i][j] > d[i][j+1] {
				d[i][j+1] = d[i][j]
			}
		}
	}
	full := 1 << m
	dp := make([]int, full)
	parV := make([]int, full)
	parI := make([]int, full)
	for i := range dp {
		dp[i] = -2
	}
	dp[0] = n
	for mask := 1; mask < full; mask++ {
		for i := 0; i < m; i++ {
			bit := 1 << i
			if mask&bit == 0 {
				continue
			}
			prev := mask ^ bit
			tmp := dp[prev]
			if tmp <= 0 {
				continue
			}
			v := d[i][tmp-1]
			if v < 0 {
				continue
			}
			if v > dp[mask] {
				dp[mask] = v
				parV[mask] = tmp
				parI[mask] = i
			}
		}
	}
	if dp[full-1] < 0 {
		fmt.Fprintln(writer, "NO")
		return
	}
	fmt.Fprintln(writer, "YES")
	sol := make([][]int, m)
	mask := full - 1
	for mask > 0 {
		i := parI[mask]
		v := parV[mask]
		tmpv := d[i][v-1]
		for j := tmpv; j < v; j++ {
			sol[i] = append(sol[i], a[j].idx)
		}
		mask ^= 1 << i
	}
	for i := 0; i < m; i++ {
		fmt.Fprint(writer, len(sol[i]))
		for _, idx := range sol[i] {
			fmt.Fprint(writer, " ", idx+1)
		}
		fmt.Fprintln(writer)
	}
}
