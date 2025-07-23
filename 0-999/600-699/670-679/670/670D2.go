package main

import (
	"bufio"
	"fmt"
	"os"
)

func canMake(a, b []int64, k, x int64) bool {
	var need int64
	for i := range a {
		required := a[i] * x
		if required > b[i] {
			need += required - b[i]
			if need > k {
				return false
			}
		}
	}
	return need <= k
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}
	var lo int64 = 0
	var hi int64 = 2000000000
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if canMake(a, b, k, mid) {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	fmt.Println(lo)
}
