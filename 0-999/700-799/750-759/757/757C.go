package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	// counts[type] -> map[gym]count
	typeCount := make(map[int]map[int]int)
	for i := 0; i < n; i++ {
		var g int
		fmt.Fscan(in, &g)
		for j := 0; j < g; j++ {
			var t int
			fmt.Fscan(in, &t)
			mp, ok := typeCount[t]
			if !ok {
				mp = make(map[int]int)
				typeCount[t] = mp
			}
			mp[i]++
		}
	}

	groups := make(map[string]int)
	buf := &bytes.Buffer{}
	idx := make([]int, 0)
	for _, mp := range typeCount {
		buf.Reset()
		idx = idx[:0]
		for gym := range mp {
			idx = append(idx, gym)
		}
		sort.Ints(idx)
		for _, gym := range idx {
			fmt.Fprintf(buf, "%d:%d|", gym, mp[gym])
		}
		key := buf.String()
		groups[key]++
	}

	absent := m - len(typeCount)
	if absent > 0 {
		groups[""] += absent
	}

	// precompute factorials up to m
	fact := make([]int64, m+1)
	fact[0] = 1
	for i := 1; i <= m; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}

	ans := int64(1)
	for _, s := range groups {
		ans = ans * fact[s] % mod
	}

	fmt.Println(ans)
}
