package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	mod1  int64 = 1000000007
	mod2  int64 = 1000000009
	base1 int64 = 911382323
	base2 int64 = 972663749
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		var s string
		fmt.Fscan(in, &s)

		zeros := make([]int, n+1)
		hash1 := make([]int64, n+1)
		hash2 := make([]int64, n+1)
		pow1 := make([]int64, n+1)
		pow2 := make([]int64, n+1)
		zeroHash1 := make([]int64, n+1)
		zeroHash2 := make([]int64, n+1)
		oneHash1 := make([]int64, n+1)
		oneHash2 := make([]int64, n+1)

		pow1[0], pow2[0] = 1, 1
		for i := 1; i <= n; i++ {
			pow1[i] = pow1[i-1] * base1 % mod1
			pow2[i] = pow2[i-1] * base2 % mod2
			zeroHash1[i] = (zeroHash1[i-1]*base1 + 1) % mod1
			zeroHash2[i] = (zeroHash2[i-1]*base2 + 1) % mod2
			oneHash1[i] = (oneHash1[i-1]*base1 + 2) % mod1
			oneHash2[i] = (oneHash2[i-1]*base2 + 2) % mod2
		}

		for i := 1; i <= n; i++ {
			val := int64(1)
			if s[i-1] == '1' {
				val = 2
			}
			hash1[i] = (hash1[i-1]*base1 + val) % mod1
			hash2[i] = (hash2[i-1]*base2 + val) % mod2
			zeros[i] = zeros[i-1]
			if s[i-1] == '0' {
				zeros[i]++
			}
		}

		uniq := make(map[[2]int64]struct{})
		for i := 0; i < m; i++ {
			var l, r int
			fmt.Fscan(in, &l, &r)
			count0 := zeros[r] - zeros[l-1]
			segLen := r - l + 1
			count1 := segLen - count0

			h1 := hash1[l-1]
			h2 := hash2[l-1]

			h1 = (h1*pow1[count0] + zeroHash1[count0]) % mod1
			h2 = (h2*pow2[count0] + zeroHash2[count0]) % mod2

			h1 = (h1*pow1[count1] + oneHash1[count1]) % mod1
			h2 = (h2*pow2[count1] + oneHash2[count1]) % mod2

			suffixLen := n - r
			sub1 := (hash1[n] - hash1[r]*pow1[suffixLen]%mod1 + mod1) % mod1
			sub2 := (hash2[n] - hash2[r]*pow2[suffixLen]%mod2 + mod2) % mod2

			h1 = (h1*pow1[suffixLen] + sub1) % mod1
			h2 = (h2*pow2[suffixLen] + sub2) % mod2

			uniq[[2]int64{h1, h2}] = struct{}{}
		}

		fmt.Fprintln(out, len(uniq))
	}
}
