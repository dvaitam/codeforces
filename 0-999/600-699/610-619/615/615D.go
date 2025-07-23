package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 1000000007

// fast exponentiation modulo mod
func modPow(base, exp, mod int64) int64 {
	res := int64(1)
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var m int
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return
	}
	counts := make(map[int]int, m)
	for i := 0; i < m; i++ {
		var p int
		fmt.Fscan(reader, &p)
		counts[p]++
	}

	type pair struct{ p, c int }
	arr := make([]pair, 0, len(counts))
	for p, c := range counts {
		arr = append(arr, pair{p, c})
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].p < arr[j].p })

	k := len(arr)
	modMinus1 := MOD - 1
	pref := make([]int64, k+1)
	suff := make([]int64, k+1)
	pref[0] = 1
	for i := 0; i < k; i++ {
		pref[i+1] = pref[i] * int64(arr[i].c+1) % modMinus1
	}
	suff[k] = 1
	for i := k - 1; i >= 0; i-- {
		suff[i] = suff[i+1] * int64(arr[i].c+1) % modMinus1
	}

	ans := int64(1)
	for i := 0; i < k; i++ {
		a := int64(arr[i].c)
		expPart := a * (a + 1) / 2 % modMinus1
		other := pref[i] * suff[i+1] % modMinus1
		exp := expPart * other % modMinus1
		ans = ans * modPow(int64(arr[i].p), exp, MOD) % MOD
	}

	fmt.Fprintln(writer, ans)
}
