package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 998244353

func powMod(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	vals := make([]int64, n)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		vals[i] = int64(x)
	}

	sort.Slice(vals, func(i, j int) bool { return vals[i] > vals[j] })

	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + vals[i-1]
	}

	invN := powMod(int64(n), MOD-2)

	results := make([]int64, n)
	for k := 1; k <= n; k++ {
		q := (n - k) / k
		r := (n - k) % k
		var sum int64
		for t := 0; t < q; t++ {
			l := k + t*k + 1
			rIndex := k + (t+1)*k
			segment := prefix[rIndex] - prefix[l-1]
			sum += int64(t+1) * segment
		}
		if r > 0 {
			l := k + q*k + 1
			rIndex := k + q*k + r
			segment := prefix[rIndex] - prefix[l-1]
			sum += int64(q+1) * segment
		}
		results[k-1] = sum % MOD * invN % MOD
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, results[i])
	}
}
