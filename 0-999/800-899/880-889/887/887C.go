package main

import (
	"bufio"
	"fmt"
	"os"
)

func solved(s []int) bool {
	for i := 0; i < 24; i += 4 {
		if s[i] != s[i+1] || s[i] != s[i+2] || s[i] != s[i+3] {
			return false
		}
	}
	return true
}

func apply(s []int, m []int) []int {
	res := make([]int, 24)
	for i, j := range m {
		res[j] = s[i]
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	s := make([]int, 24)
	for i := 0; i < 24; i++ {
		if _, err := fmt.Fscan(in, &s[i]); err != nil {
			return
		}
	}

	rotations := [][]int{
		{1, 3, 0, 2, 17, 16, 6, 7, 4, 5, 10, 11, 8, 9, 14, 15, 13, 12, 18, 19, 20, 21, 22, 23}, // U
		{0, 1, 2, 3, 4, 5, 10, 11, 8, 9, 14, 15, 12, 13, 19, 18, 16, 17, 7, 6, 21, 23, 20, 22}, // D
		{0, 1, 12, 14, 4, 3, 6, 2, 9, 11, 8, 10, 21, 13, 20, 15, 16, 17, 18, 19, 5, 7, 22, 23}, // F
		{6, 4, 2, 3, 22, 5, 23, 7, 8, 9, 10, 11, 12, 0, 14, 1, 18, 16, 19, 17, 20, 21, 15, 13}, // B
		{8, 1, 10, 3, 5, 7, 4, 6, 20, 9, 22, 11, 12, 13, 14, 15, 2, 17, 0, 19, 18, 21, 16, 23}, // L
		{0, 19, 2, 17, 4, 5, 6, 7, 8, 1, 10, 3, 13, 15, 12, 14, 16, 23, 18, 21, 20, 9, 22, 11}, // R
	}

	for _, m := range rotations {
		t := apply(s, m)
		if solved(t) {
			fmt.Println("YES")
			return
		}
		t = apply(apply(apply(s, m), m), m)
		if solved(t) {
			fmt.Println("YES")
			return
		}
	}
	fmt.Println("NO")
}
