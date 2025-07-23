package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m)

	if m == 1 {
		if n == 1 {
			fmt.Println(1)
		} else {
			fmt.Println(2)
		}
		return
	}
	if m == n {
		fmt.Println(n - 1)
		return
	}

	left := m - 1
	right := n - m
	if left >= right {
		fmt.Println(m - 1)
	} else {
		fmt.Println(m + 1)
	}
}
