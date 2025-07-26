package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	parent := make([]int, n+1)
	for i := 2; i <= n; i++ {
		fmt.Fscan(in, &parent[i])
	}
	children := make([]int, n+1)
	for i := 2; i <= n; i++ {
		children[parent[i]]++
	}
	leaves := []int{}
	for i := 1; i <= n; i++ {
		if children[i] == 0 {
			leaves = append(leaves, i)
		}
	}
	leafIndex := make(map[int]int)
	for idx, v := range leaves {
		leafIndex[v] = idx
	}
	m := len(leaves)

	hunger := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &hunger[i])
	}

	leafH := make([]int64, m)
	for idx, v := range leaves {
		leafH[idx] = hunger[v]
	}

	var q int
	fmt.Fscan(in, &q)

	calc := func() int64 {
		vals := make([]int64, m)
		copy(vals, leafH)
		sort.Slice(vals, func(i, j int) bool { return vals[i] > vals[j] })
		var ans int64
		for i, v := range vals {
			if v == 0 {
				continue
			}
			t := int64(i+1) + (v-1)*int64(m)
			if t > ans {
				ans = t
			}
		}
		return ans % MOD
	}

	fmt.Fprintln(out, calc())
	for ; q > 0; q-- {
		var v int
		var x int64
		fmt.Fscan(in, &v, &x)
		if idx, ok := leafIndex[v]; ok {
			leafH[idx] = x
		}
		fmt.Fprintln(out, calc())
	}
}
