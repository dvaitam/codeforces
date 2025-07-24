package main

import (
	"bufio"
	"fmt"
	"os"
)

func primesUpTo(n int) []int {
	primes := []int{}
	for i := 2; i <= n; i++ {
		isPrime := true
		for j := 2; j*j <= i; j++ {
			if i%j == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			primes = append(primes, i)
		}
	}
	return primes
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	primes := primesUpTo(100)

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := range a {
			fmt.Fscan(in, &a[i])
		}

		seen := make(map[int64]bool)
		dup := false
		for _, v := range a {
			if seen[v] {
				dup = true
				break
			}
			seen[v] = true
		}
		if dup {
			fmt.Fprintln(out, "NO")
			continue
		}

		ok := true
		for _, p := range primes {
			if p > n {
				break
			}
			cnt := make([]int, p)
			for _, v := range a {
				cnt[int(v%int64(p))]++
			}
			forb := make([]bool, p)
			for r, c := range cnt {
				if c >= 2 {
					forb[(p-r)%p] = true
				}
			}
			all := true
			for _, b := range forb {
				if !b {
					all = false
					break
				}
			}
			if all {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
