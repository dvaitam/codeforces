package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxN = 100000

func sieve(limit int) []bool {
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}
	return isPrime
}

func lastPrimeAtMost(x int, isPrime []bool) int {
	if x < 2 {
		return 2
	}
	for x >= 2 && !isPrime[x] {
		x--
	}
	return x
}

func buildPermutation(n int, isPrime []bool) []int {
	need := n/3 - 1
	if need < 0 {
		need = 0
	}

	perm := make([]int, 0, n)
	used := make([]bool, n+1)

	if need > 0 {
		center := lastPrimeAtMost(n/2, isPrime)
		start := center - need/2
		end := start + need - 1
		if start < 1 {
			start = 1
			end = start + need - 1
		}
		if end > n {
			end = n
			start = end - need + 1
		}

		mid := need / 2
		for i := 0; i < need; i++ {
			idx := mid
			if i > 0 {
				if i%2 == 1 {
					idx = mid - (i+1)/2
				} else {
					idx = mid + i/2
				}
			}
			val := start + idx
			perm = append(perm, val)
			used[val] = true
		}
	}

	for i := 1; i <= n; i++ {
		if !used[i] {
			perm = append(perm, i)
		}
	}

	return perm
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	isPrime := sieve(maxN)

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		perm := buildPermutation(n, isPrime)
		for i, v := range perm {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
