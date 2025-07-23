package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	mod1 int64 = 1000000007
	mod2 int64 = 1000000009
	base int64 = 911382323
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)
	m := len(s)
	// precompute powers and prefix hashes
	pow1 := make([]int64, m+1)
	pow2 := make([]int64, m+1)
	h1 := make([]int64, m+1)
	h2 := make([]int64, m+1)
	pow1[0], pow2[0] = 1, 1
	for i := 1; i <= m; i++ {
		ch := int64(s[i-1] - 'a' + 1)
		pow1[i] = pow1[i-1] * base % mod1
		pow2[i] = pow2[i-1] * base % mod2
		h1[i] = (h1[i-1]*base + ch) % mod1
		h2[i] = (h2[i-1]*base + ch) % mod2
	}
	getHash := func(l, r int) (int64, int64) {
		x1 := (h1[r] - h1[l]*pow1[r-l]) % mod1
		if x1 < 0 {
			x1 += mod1
		}
		x2 := (h2[r] - h2[l]*pow2[r-l]) % mod2
		if x2 < 0 {
			x2 += mod2
		}
		return x1, x2
	}
	res := make([]int, (n+1)/2)
	for k := 1; k <= (n+1)/2; k++ {
		start := k - 1
		end := n - k
		length := end - start + 1
		best := -1
		for l := length - 1; l >= 1; l-- {
			if l%2 == 0 {
				continue
			}
			h1a, h2a := getHash(start, start+l)
			h1b, h2b := getHash(end-l+1, end+1)
			if h1a == h1b && h2a == h2b {
				best = l
				break
			}
		}
		res[k-1] = best
	}
	for i, v := range res {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(v)
	}
	fmt.Println()
}
