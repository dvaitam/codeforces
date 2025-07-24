package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// calcCost computes minimal number of type1 spells needed for given prefix array.
// This is a naive O(n^2) implementation derived from the problem statement.
func calcCost(arr []int) int {
	b := append([]int(nil), arr...)
	sort.Ints(b)
	n := len(b)
	best := int(1 << 60)
	for L := 0; L <= n+1; L++ {
		idx := 0
		costPre := 0
		feasible := true
		for v := 1; v < L; v++ {
			for idx < n && b[idx] < v {
				idx++
			}
			if idx == n {
				feasible = false
				break
			}
			costPre += b[idx] - v
			idx++
		}
		if !feasible {
			break
		}
		cost := costPre
		for j := idx; j < n; j++ {
			if b[j] == L {
				cost += 1
			} else if b[j] > L {
				cost += b[j] - L
			}
		}
		if cost < best {
			best = cost
		}
	}
	return best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		prefix := make([]int, 0, n)
		for i := 0; i < n; i++ {
			prefix = append(prefix, a[i])
			ans := calcCost(prefix)
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans)
		}
		fmt.Fprintln(out)
	}
}
