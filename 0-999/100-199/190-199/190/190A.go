package main

import "fmt"

func main() {
	var n, m int
	if _, err := fmt.Scan(&n, &m); err != nil {
		return
	}
	// Impossible if there are children but no grown-ups
	if n == 0 && m > 0 {
		fmt.Println("Impossible")
		return
	}
	// Minimum fare: maximize free rides (one free child per grown-up)
	minFare := n
	if m > n {
		minFare = m
	}
	// Maximum fare: minimize free rides by grouping children
	var maxFare int
	if m == 0 {
		maxFare = n
	} else {
		maxFare = n + m - 1
	}
	fmt.Printf("%d %d", minFare, maxFare)
}
