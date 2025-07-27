package main

import (
	"bufio"
	"fmt"
	"os"
)

func comb2(n int64) int64 {
	if n < 2 {
		return 0
	}
	return n * (n - 1) / 2
}

func comb3(n int64) int64 {
	if n < 3 {
		return 0
	}
	return n * (n - 1) * (n - 2) / 6
}

func probability(n, k int64) float64 {
	d := float64(comb3(n))
	p0 := float64(comb3(n - k))
	p1 := float64(k * comb2(n-k))
	return 1 - p0/d - p1/(2*d)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int64
	var p float64
	fmt.Fscan(reader, &n, &p)
	for k := int64(0); k <= n; k++ {
		if probability(n, k) >= p-1e-12 {
			fmt.Println(k)
			return
		}
	}
}
