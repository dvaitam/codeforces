package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemC.txt.
// We form a mapping from each value a[i] to b[i]. The mapping
// consists of disjoint cycles. For each cycle, if no position
// in that cycle is fixed (d[i] != 0) and all a[i] != b[i], then
// the cycle can be oriented in two ways. Otherwise its orientation
// is forced. The answer is 2^(number of unfixed cycles) modulo 1e9+7.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const mod int64 = 1_000_000_007
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		b := make([]int, n)
		d := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &d[i])
		}

		pos := make([]int, n+1)
		for i := 0; i < n; i++ {
			pos[a[i]] = i
		}

		visited := make([]bool, n+1)
		ans := int64(1)
		for x := 1; x <= n; x++ {
			if visited[x] {
				continue
			}
			cur := x
			indices := make([]int, 0)
			for !visited[cur] {
				visited[cur] = true
				idx := pos[cur]
				indices = append(indices, idx)
				cur = b[idx]
			}
			forced := false
			for _, idx := range indices {
				if d[idx] != 0 || a[idx] == b[idx] {
					forced = true
					break
				}
			}
			if !forced {
				ans = (ans * 2) % mod
			}
		}
		fmt.Fprintln(out, ans)
	}
}
