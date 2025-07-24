package main

import "fmt"

func main() {
	var n, k int64
	if _, err := fmt.Scan(&n, &k); err != nil {
		return
	}

	if k == 0 || k == n {
		fmt.Printf("0 0\n")
		return
	}

	minGood := int64(1)
	maxGood := n - k
	if maxGood > 2*k {
		maxGood = 2 * k
	}

	fmt.Printf("%d %d\n", minGood, maxGood)
}
