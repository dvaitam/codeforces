package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

var (
	n, k int
	a    []int
	cost []int
	memo map[state]int
)

type state struct {
	idx int
	s1  uint64
	s2  uint64
}

func dfs(idx int, s1, s2 uint64) int {
	if idx == n {
		return 0
	}
	st := state{idx, s1, s2}
	if v, ok := memo[st]; ok {
		return v
	}
	book := a[idx]
	var mask1, mask2 uint64
	if book <= 64 {
		mask1 = 1 << (book - 1)
	} else {
		mask2 = 1 << (book - 65)
	}
	// Check if the book is already in the set
	if (book <= 64 && (s1&mask1) != 0) || (book > 64 && (s2&mask2) != 0) {
		res := dfs(idx+1, s1, s2)
		memo[st] = res
		return res
	}
	size := bits.OnesCount64(s1) + bits.OnesCount64(s2)
	best := int(^uint(0) >> 1) // Max int
	if int(size) < k {
		ns1, ns2 := s1, s2
		if book <= 64 {
			ns1 |= mask1
		} else {
			ns2 |= mask2
		}
		val := cost[book] + dfs(idx+1, ns1, ns2)
		if val < best {
			best = val
		}
	} else {
		// Try replacing each existing book
		for b := s1; b != 0; b &= b - 1 {
			i := bits.TrailingZeros64(b)
			ns1 := (s1 &^ (1 << uint(i)))
			ns2 := s2
			if book <= 64 {
				ns1 |= mask1
			} else {
				ns2 |= mask2
			}
			val := cost[book] + dfs(idx+1, ns1, ns2)
			if val < best {
				best = val
			}
		}
		for b := s2; b != 0; b &= b - 1 {
			i := bits.TrailingZeros64(b)
			ns1 := s1
			ns2 := (s2 &^ (1 << uint(i)))
			if book <= 64 {
				ns1 |= mask1
			} else {
				ns2 |= mask2
			}
			val := cost[book] + dfs(idx+1, ns1, ns2)
			if val < best {
				best = val
			}
		}
	}
	memo[st] = best
	return best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	cost = make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &cost[i])
	}
	memo = make(map[state]int)
	res := dfs(0, 0, 0)
	fmt.Println(res)
}
