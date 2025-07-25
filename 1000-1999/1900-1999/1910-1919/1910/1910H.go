package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	const K = 9
	pow10 := make([]int, K+1)
	pow10[0] = 1
	for i := 1; i <= K; i++ {
		pow10[i] = pow10[i-1] * 10
	}

	L := make([]int64, K+1)
	rem := make([][]int, K+1)
	for i := 1; i <= K; i++ {
		rem[i] = make([]int, n)
	}
	var sumA int64
	for idx, v := range a {
		sumA += int64(v)
		for k := 1; k <= K; k++ {
			L[k] += int64(v / pow10[k])
			rem[k][idx] = v % pow10[k]
		}
	}
	for k := 1; k <= K; k++ {
		sort.Ints(rem[k])
	}

	for _, v := range a {
		res := int64(n)*int64(v) + sumA
		for k := 1; k <= K; k++ {
			p := pow10[k]
			r := v % p
			cnt := int64(len(rem[k]) - sort.SearchInts(rem[k], p-r))
			res -= 9 * (int64(n*(v/p)) + L[k] + cnt)
		}
		fmt.Fprint(out, res, " ")
	}
	fmt.Fprintln(out)
}
