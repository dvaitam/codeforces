package main

import (
	"fmt"
)

func main() {
	var d, sumTime int
	if _, err := fmt.Scan(&d, &sumTime); err != nil {
		return
	}
	mins := make([]int, d)
	maxs := make([]int, d)
	for i := 0; i < d; i++ {
		fmt.Scan(&mins[i], &maxs[i])
	}
	// check feasibility
	sumMin := 0
	sumMax := 0
	for i := 0; i < d; i++ {
		sumMin += mins[i]
		sumMax += maxs[i]
	}
	if sumTime < sumMin || sumTime > sumMax {
		fmt.Println("NO")
		return
	}
	// construct a valid allocation
	rem := sumTime - sumMin
	res := make([]int, d)
	for i := 0; i < d; i++ {
		extra := maxs[i] - mins[i]
		if extra > rem {
			extra = rem
		}
		res[i] = mins[i] + extra
		rem -= extra
	}
	fmt.Println("YES")
	for i, v := range res {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(v)
	}
	fmt.Println()
}
