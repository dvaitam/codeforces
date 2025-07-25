package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func countPairs(arr []int, c int) int64 {
	var res int64
	for i := 0; i < len(arr); i++ {
		limit := 2*c - arr[i]
		j := sort.Search(i+1, func(k int) bool { return arr[k] > limit }) - 1
		if j > i {
			j = i
		}
		if j >= 0 {
			res += int64(j + 1)
		}
	}
	return res
}

func solveCase(n, c int, s []int) int64 {
	total := int64(c+1) * int64(c+2) / 2

	var aCount int64
	for _, v := range s {
		lo := v - c
		if lo < 0 {
			lo = 0
		}
		hi := v / 2
		if hi > c {
			hi = c
		}
		if lo <= hi {
			aCount += int64(hi - lo + 1)
		}
	}

	var bCount int64
	for _, v := range s {
		bCount += int64(c + 1 - v)
	}

	var even, odd []int
	for _, v := range s {
		if v%2 == 0 {
			even = append(even, v)
		} else {
			odd = append(odd, v)
		}
	}

	var abCount int64
	abCount += countPairs(even, c)
	abCount += countPairs(odd, c)

	return total - aCount - bCount + abCount
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, c int
		fmt.Fscan(in, &n, &c)
		s := make([]int, n)
		for i := range s {
			fmt.Fscan(in, &s[i])
		}
		ans := solveCase(n, c, s)
		fmt.Fprintln(out, ans)
	}
}
