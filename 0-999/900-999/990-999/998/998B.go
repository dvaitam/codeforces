package main

import (
	"fmt"
	"sort"
)

func main() {
	var n, B int
	if _, err := fmt.Scan(&n, &B); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&a[i])
	}
	var costs []int
	odd, even := 0, 0
	for i := 0; i < n-1; i++ {
		if a[i]%2 == 0 {
			even++
		} else {
			odd++
		}
		if odd == even {
			diff := a[i] - a[i+1]
			if diff < 0 {
				diff = -diff
			}
			costs = append(costs, diff)
		}
	}
	sort.Ints(costs)
	total := 0
	cnt := 0
	for _, c := range costs {
		if total+c <= B {
			total += c
			cnt++
		} else {
			break
		}
	}
	fmt.Println(cnt)
}
