package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	pos   int64
	count int64
}

func primeFactors(x int64) []int64 {
	factors := make([]int64, 0)
	for i := int64(2); i*i <= x; i++ {
		if x%i == 0 {
			factors = append(factors, i)
			for x%i == 0 {
				x /= i
			}
		}
	}
	if x > 1 {
		factors = append(factors, x)
	}
	return factors
}

func groupCost(group []pair, k int64) int64 {
	mid := (k + 1) / 2
	var median int64
	var acc int64
	for _, p := range group {
		acc += p.count
		if acc >= mid {
			median = p.pos
			break
		}
	}
	var cost int64
	for _, p := range group {
		if p.pos > median {
			cost += (p.pos - median) * p.count
		} else {
			cost += (median - p.pos) * p.count
		}
	}
	return cost
}

func costForK(a []int64, k int64) int64 {
	var cost int64
	var group []pair
	var sum int64
	for i, v := range a {
		val := v % k
		for val > 0 {
			need := k - sum
			take := val
			if take > need {
				take = need
			}
			group = append(group, pair{pos: int64(i), count: take})
			sum += take
			val -= take
			if sum == k {
				cost += groupCost(group, k)
				group = group[:0]
				sum = 0
			}
		}
	}
	return cost
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int64, n)
	var total int64
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		total += a[i]
	}
	if total <= 1 {
		fmt.Fprintln(writer, -1)
		return
	}
	factors := primeFactors(total)
	best := int64(1<<63 - 1)
	for _, k := range factors {
		c := costForK(a, k)
		if c < best {
			best = c
		}
	}
	fmt.Fprintln(writer, best)
}
