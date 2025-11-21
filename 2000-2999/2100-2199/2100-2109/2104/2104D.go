package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const maxN = 400000
const sieveLimit = 7000000

var primePrefix [maxN + 1]int64

func init() {
	primes := make([]int, 0, maxN)
	isComposite := make([]bool, sieveLimit+1)
	for i := 2; i <= sieveLimit && len(primes) < maxN; i++ {
		if !isComposite[i] {
			primes = append(primes, i)
			if int64(i)*int64(i) <= sieveLimit {
				for j := i * i; j <= sieveLimit; j += i {
					isComposite[j] = true
				}
			}
		}
	}
	if len(primes) < maxN {
		panic("not enough primes computed")
	}
	primePrefix[0] = 0
	for i := 1; i <= maxN; i++ {
		primePrefix[i] = primePrefix[i-1] + int64(primes[i-1])
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
		prefix := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			prefix[i] = prefix[i-1] + a[i-1]
		}
		best := 0
		for l := n; l >= 0; l-- {
			if primePrefix[l] <= prefix[l] {
				best = l
				break
			}
		}
		fmt.Fprintln(out, n-best)
	}
}
