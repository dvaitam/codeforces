package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func combParity(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	if (k & (n - k)) == 0 {
		return 1
	}
	return 0
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var k, m int
		fmt.Fscan(in, &k, &m)
		var s string
		fmt.Fscan(in, &s)
		if len(s) != k {
			k = len(s)
		}

		var nVal big.Int
		nVal.SetString(s, 2)
		var H big.Int
		H.Set(&nVal)
		H.Add(&H, big.NewInt(1))
		H.Rsh(&H, 1)

		prefix := big.NewInt(0)
		t0 := -1
		totalBits := 0
		for i := 0; i < k; i++ {
			if s[i] == '1' {
				totalBits++
				exp := uint(k - 1 - i)
				var bit big.Int
				bit.Lsh(big.NewInt(1), exp)
				prefix.Add(prefix, &bit)
				if t0 == -1 && prefix.Cmp(&H) >= 0 {
					t0 = totalBits
				}
			}
		}

		if t0 == -1 {
			t0 = totalBits
		}

		ans := 0
		for i := 0; i < m; i++ {
			A := combParity(t0+i, i)
			B := combParity(t0+i-1, i-1)
			prefixParity := A ^ B

			L := totalBits - t0
			M := m - 1 - i
			suffixParity := combParity(L+M, M)

			if prefixParity == 1 && suffixParity == 1 {
				ans ^= i
			}
		}
		fmt.Fprintln(out, ans)
	}
}
