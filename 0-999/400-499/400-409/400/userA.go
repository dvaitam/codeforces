package main

import "fmt"

func main() {
	var t int
	fmt.Scan(&t)
	for i := 0; i < t; i++ {
		var s string
		fmt.Scan(&s)
		var pairs [][]int
		for _, a := range []int{1, 2, 3, 4, 6, 12} {
			b := 12 / a
			if check(s, a, b) {
				pairs = append(pairs, []int{a, b})
			}
		}
		fmt.Print(len(pairs))
		for _, p := range pairs {
			fmt.Printf(" %dx%d", p[0], p[1])
		}
		fmt.Println()
	}
}

func check(s string, a, b int) bool {
	for j := 0; j < b; j++ {
		allX := true
		for k := 0; k < a; k++ {
			pos := k*b + j
			if s[pos] != 'X' {
				allX = false
				break
			}
		}
		if allX {
			return true
		}
	}
	return false
}