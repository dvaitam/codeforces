package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt for contest 595A.
// It counts how many flats have at least one window lit.
func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	count := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			var a, b int
			fmt.Fscan(in, &a, &b)
			if a == 1 || b == 1 {
				count++
			}
		}
	}
	fmt.Println(count)
}
