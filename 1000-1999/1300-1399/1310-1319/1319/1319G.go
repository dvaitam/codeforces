package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const (
	mod  = int64(1000000007)
	base = int64(911382323)
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var t string
	fmt.Fscan(in, &t)

	prefixZero := make([]int, n+1)
	prefixOne := make([]int, n+1)
	prefixParity := make([]int, n+1) // ones parity
	zeroPos := make([]int, 0)
	zeroParity := make([]int, 0)

	for i := 1; i <= n; i++ {
		if t[i-1] == '0' {
			prefixZero[i] = prefixZero[i-1] + 1
			prefixOne[i] = prefixOne[i-1]
			zeroPos = append(zeroPos, i)
			p := prefixOne[i] % 2
			zeroParity = append(zeroParity, p)
		} else {
			prefixZero[i] = prefixZero[i-1]
			prefixOne[i] = prefixOne[i-1] + 1
		}
		prefixParity[i] = prefixOne[i] % 2
	}

	m := len(zeroPos)
	powBase := make([]int64, m+1)
	powSum := make([]int64, m+1)
	powBase[0] = 1
	powSum[0] = 0
	for i := 1; i <= m; i++ {
		powBase[i] = powBase[i-1] * base % mod
		powSum[i] = (powSum[i-1]*base + 1) % mod
	}

	prefixHash := make([]int64, m+1)
	for i := 0; i < m; i++ {
		prefixHash[i+1] = (prefixHash[i]*base + int64(zeroParity[i])) % mod
	}

	var q int
	fmt.Fscan(in, &q)

	for ; q > 0; q-- {
		var l1, l2, length int
		fmt.Fscan(in, &l1, &l2, &length)
		r1 := l1 + length - 1
		r2 := l2 + length - 1

		z1 := prefixZero[r1] - prefixZero[l1-1]
		z2 := prefixZero[r2] - prefixZero[l2-1]
		if z1 != z2 {
			fmt.Fprintln(out, "NO")
			continue
		}
		z := z1
		if z == 0 {
			fmt.Fprintln(out, "YES")
			continue
		}

		// indices of zeros in substring
		startIdx1 := sort.SearchInts(zeroPos, l1)
		endIdx1 := sort.SearchInts(zeroPos, r1+1) - 1
		startIdx2 := sort.SearchInts(zeroPos, l2)
		endIdx2 := sort.SearchInts(zeroPos, r2+1) - 1

		hash1 := getHash(prefixHash, powBase, startIdx1, endIdx1)
		hash2 := getHash(prefixHash, powBase, startIdx2, endIdx2)

		lenZeros := endIdx1 - startIdx1 + 1 // same as z

		if prefixParity[l1-1]%2 == 1 {
			hash1 = (powSum[lenZeros] - hash1 + mod) % mod
		}
		if prefixParity[l2-1]%2 == 1 {
			hash2 = (powSum[lenZeros] - hash2 + mod) % mod
		}

		if hash1 == hash2 {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

func getHash(prefixHash, powBase []int64, l, r int) int64 {
	if l > r {
		return 0
	}
	res := (prefixHash[r+1] - prefixHash[l]*powBase[r-l+1]) % mod
	if res < 0 {
		res += mod
	}
	return res
}
