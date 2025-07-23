package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	a := make([]int64, n)
	b := make([]int64, n)
	minA := int64(1<<63 - 1)
	maxB := int64(0)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] < minA {
			minA = a[i]
		}
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
		if b[i] > maxB {
			maxB = b[i]
		}
	}
	// upper bound for answer
	hi := (maxB + k) / minA
	hi += k // safety margin
	lo := int64(0)
	for lo < hi {
		mid := (lo + hi + 1) / 2
		need := int64(0)
		for i := 0; i < n; i++ {
			req := a[i] * mid
			if req > b[i] {
				need += req - b[i]
				if need > k {
					break
				}
			}
		}
		if need <= k {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	fmt.Println(lo)
}
