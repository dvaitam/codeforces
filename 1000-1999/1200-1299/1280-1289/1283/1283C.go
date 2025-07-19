package main

import (
	"fmt"
)

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	v := make([]int, n)
	used := make([]bool, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&v[i])
		if v[i] != 0 {
			used[v[i]-1] = true
		}
	}
	var givers []int
	var missing []int
	for i := 0; i < n; i++ {
		if v[i] == 0 {
			givers = append(givers, i)
		}
		if !used[i] {
			missing = append(missing, i)
		}
	}
	k := len(givers)
	for i := 0; i < k; i++ {
		j := (i + 1) % k
		v[givers[i]] = missing[j] + 1
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(v[i])
	}
	fmt.Println()
}
