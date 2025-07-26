package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	// precompute divisors for each index
	divisors := make([][]int, n+1)
	for d := 1; d <= n; d++ {
		for m := d; m <= n; m += d {
			divisors[m] = append(divisors[m], d)
		}
	}

	// precompute powers of 2 modulo mod
	pow2 := make([]int64, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = (pow2[i-1] * 2) % mod
	}

	// sort indices by value descending
	idxs := make([]int, n)
	for i := 0; i < n; i++ {
		idxs[i] = i + 1
	}
	sort.Slice(idxs, func(i, j int) bool {
		if a[idxs[i]] == a[idxs[j]] {
			return idxs[i] < idxs[j]
		}
		return a[idxs[i]] > a[idxs[j]]
	})

	banned := make([]bool, n+1)
	avail := n
	ans := int64(0)

	visited := make([]bool, n+1)
	used := make([]int, 0)

	for i := 0; i < n; {
		v := a[idxs[i]]
		group := make([]int, 0)
		for i < n && a[idxs[i]] == v {
			group = append(group, idxs[i])
			i++
		}

		// gather unique available divisors for this group
		k := 0
		used = used[:0]
		for _, x := range group {
			for _, d := range divisors[x] {
				if !banned[d] && !visited[d] {
					visited[d] = true
					used = append(used, d)
					k++
				}
			}
		}

		if k > 0 {
			cnt := pow2[avail] - pow2[avail-k]
			if cnt < 0 {
				cnt += mod
			}
			ans = (ans + int64(v)*cnt) % mod
		}

		for _, d := range used {
			visited[d] = false
		}

		for _, x := range group {
			for _, d := range divisors[x] {
				if !banned[d] {
					banned[d] = true
					avail--
				}
			}
		}
	}

	fmt.Fprintln(out, ans%mod)
}
