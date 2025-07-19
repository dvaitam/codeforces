package main

import (
	"fmt"
	"strings"
)

func main() {
	var s1, s2 string
	if _, err := fmt.Scan(&s1); err != nil {
		return
	}
	fmt.Scan(&s2)

	parts := strings.Split(s1, "|")
	if len(parts) != 2 {
		fmt.Println("Impossible")
		return
	}
	left := parts[0]
	right := parts[1]
	a, b, c := len(left), len(right), len(s2)
	// Check feasibility: total must be even and difference <= remaining weights
	if (a+b+c)%2 != 0 || abs(a-b) > c {
		fmt.Println("Impossible")
		return
	}
	// Balance the initial difference
	diff := abs(a - b)
	idx := 0
	if a < b {
		left += s2[idx : idx+diff]
	} else {
		right += s2[idx : idx+diff]
	}
	idx += diff
	// Distribute remaining weights evenly
	rem := s2[idx:]
	half := len(rem) / 2
	left += rem[:half]
	right += rem[half:]
	// Output result
	fmt.Println(left + "|" + right)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
