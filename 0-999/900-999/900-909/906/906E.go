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

	var s, t string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	n := len(s)
	if n != len(t) {
		fmt.Fprintln(out, -1)
		return
	}

	const mod1 int64 = 1000000007
	const mod2 int64 = 1000000009
	const base int64 = 911382323

	pow1 := make([]int64, n+1)
	pow2 := make([]int64, n+1)
	pow1[0], pow2[0] = 1, 1
	for i := 1; i <= n; i++ {
		pow1[i] = pow1[i-1] * base % mod1
		pow2[i] = pow2[i-1] * base % mod2
	}

	hs1 := make([]int64, n+1)
	hs2 := make([]int64, n+1)
	htRev1 := make([]int64, n+1)
	htRev2 := make([]int64, n+1)

	tr := []byte(t)
	for l, r := 0, n-1; l < r; l, r = l+1, r-1 {
		tr[l], tr[r] = tr[r], tr[l]
	}

	for i := 0; i < n; i++ {
		c := int64(s[i])
		hs1[i+1] = (hs1[i]*base + c) % mod1
		hs2[i+1] = (hs2[i]*base + c) % mod2
		cr := int64(tr[i])
		htRev1[i+1] = (htRev1[i]*base + cr) % mod1
		htRev2[i+1] = (htRev2[i]*base + cr) % mod2
	}

	getHash := func(h []int64, pow []int64, mod int64, l, r int) int64 {
		val := h[r] - h[l]*pow[r-l]%mod
		if val < 0 {
			val += mod
		}
		return val
	}

	isEqual := func(l, r int) bool {
		h1s := getHash(hs1, pow1, mod1, l, r)
		h2s := getHash(hs2, pow2, mod2, l, r)
		revL := n - r
		revR := n - l
		h1t := getHash(htRev1, pow1, mod1, revL, revR)
		h2t := getHash(htRev2, pow2, mod2, revL, revR)
		return h1s == h1t && h2s == h2t
	}

	type pair struct{ l, r int }
	res := make([]pair, 0)
	for i, r := 0, 0; i < n; {
		if s[i] == t[i] {
			i++
			if r < i {
				r = i
			}
			continue
		}
		if r < i+1 {
			r = i + 1
		}
		for r <= n && !isEqual(i, r) {
			r++
		}
		if r > n {
			fmt.Fprintln(out, -1)
			return
		}
		res = append(res, pair{i + 1, r})
		i = r
		r = i
	}

	fmt.Fprintln(out, len(res))
	for _, p := range res {
		fmt.Fprintln(out, p.l, p.r)
	}
}
