package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	A := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i])
	}
	X := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &X[i])
	}

	sort.Ints(A)

	// good[k] == true iff prefixes above bit k are pairwise distinct in A.
	var good [30]bool
	for k := 0; k < 30; k++ {
		good[k] = true
		prev := -1
		for _, a := range A {
			cur := a >> (k + 1)
			if cur == prev {
				good[k] = false
				break
			}
			prev = cur
		}
	}

	// Count pairs with xor == 0 (identical elements).
	freq := make(map[int]int, m)
	var ans int64
	for _, x := range X {
		freq[x]++
	}
	for _, v := range freq {
		ans += int64(v*(v-1)) / 2
	}

	// For each k, count pairs whose xor has msb == k.
	for k := 0; k < 30; k++ {
		if !good[k] {
			continue
		}
		mp := make(map[int][2]int)
		for _, x := range X {
			prefix := x >> (k + 1)
			bit := (x >> k) & 1
			val := mp[prefix]
			val[bit]++
			mp[prefix] = val
		}
		for _, v := range mp {
			ans += int64(v[0] * v[1])
		}
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, ans)
	out.Flush()
}
