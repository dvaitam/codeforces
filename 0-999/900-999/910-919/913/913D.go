package main

import (
	"bufio"
	"fmt"
	"os"
)

type Problem struct {
	a   int
	t   int
	idx int
}

const maxT = 10000

var problems []Problem
var n int
var Tlimit int

func check(k int) (bool, []int) {
	if k == 0 {
		return true, []int{}
	}
	buckets := make([][]int, maxT+1)
	for _, p := range problems {
		if p.a >= k {
			buckets[p.t] = append(buckets[p.t], p.idx)
		}
	}
	if countAvailable := func() int {
		c := 0
		for i := 1; i <= maxT; i++ {
			c += len(buckets[i])
		}
		return c
	}(); countAvailable < k {
		return false, nil
	}
	res := make([]int, 0, k)
	timeSum := 0
	for t := 1; t <= maxT && len(res) < k; t++ {
		for _, idx := range buckets[t] {
			res = append(res, idx)
			timeSum += t
			if len(res) == k {
				break
			}
		}
	}
	if len(res) < k || timeSum > Tlimit {
		return false, nil
	}
	return true, res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &Tlimit)
	problems = make([]Problem, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &problems[i].a, &problems[i].t)
		problems[i].idx = i + 1
	}

	lo, hi := 0, n
	best := 0
	var bestSet []int
	for lo <= hi {
		mid := (lo + hi) / 2
		ok, set := check(mid)
		if ok {
			best = mid
			bestSet = set
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}

	fmt.Println(best)
	fmt.Println(best)
	for i, id := range bestSet {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(id)
	}
	if best > 0 {
		fmt.Println()
	}
}
