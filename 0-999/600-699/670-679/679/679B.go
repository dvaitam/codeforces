package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Pair struct {
	cnt int64
	sum int64
}

var memo map[int64]Pair

func cubeRoot(x int64) int64 {
	r := int64(math.Cbrt(float64(x)))
	for (r+1)*(r+1)*(r+1) <= x {
		r++
	}
	for r*r*r > x {
		r--
	}
	return r
}

func best(x int64) Pair {
	if val, ok := memo[x]; ok {
		return val
	}
	if x == 0 {
		return Pair{0, 0}
	}
	t := cubeRoot(x)
	res1 := best(t*t*t - 1)
	res2 := best(x - t*t*t)
	res2.cnt++
	res2.sum += t * t * t
	var res Pair
	if res1.cnt > res2.cnt || (res1.cnt == res2.cnt && res1.sum > res2.sum) {
		res = res1
	} else {
		res = res2
	}
	memo[x] = res
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var m int64
	fmt.Fscan(in, &m)
	memo = make(map[int64]Pair)
	ans := best(m)
	fmt.Printf("%d %d", ans.cnt, ans.sum)
}
