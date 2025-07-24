package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func mobius(n int) []int {
	mu := make([]int, n+1)
	if n >= 1 {
		mu[1] = 1
	}
	primes := make([]int, 0)
	isComp := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		if !isComp[i] {
			primes = append(primes, i)
			mu[i] = -1
		}
		for _, p := range primes {
			if i*p > n {
				break
			}
			isComp[i*p] = true
			if i%p == 0 {
				mu[i*p] = 0
				break
			}
			mu[i*p] = -mu[i]
		}
	}
	return mu
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	maxA := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		if a[i] > maxA {
			maxA = a[i]
		}
	}
	sort.Ints(a)

	pos := make([]int, maxA+1)
	for i := range pos {
		pos[i] = -1
	}
	for i, v := range a {
		pos[v] = i
	}

	mu := mobius(maxA)
	pair := make([]int64, maxA+1)
	for d := 1; d <= maxA; d++ {
		prefix := 0
		cnt := 0
		var sum int64
		for v := d; v <= maxA; v += d {
			idx := pos[v]
			if idx != -1 {
				sum += int64(cnt*idx - prefix - cnt)
				prefix += idx
				cnt++
			}
		}
		pair[d] = sum
	}

	var ans int64
	for d := 1; d <= maxA; d++ {
		if mu[d] == 0 {
			continue
		}
		ans += int64(mu[d]) * pair[d]
	}
	fmt.Fprintln(out, ans)
}
