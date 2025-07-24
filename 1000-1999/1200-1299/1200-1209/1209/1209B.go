package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	var s string
	fmt.Fscan(reader, &s)
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i], &b[i])
	}

	state := make([]int, n)
	cur := 0
	for i := 0; i < n; i++ {
		if s[i] == '1' {
			state[i] = 1
			cur++
		}
	}

	maxOn := cur
	for t := 1; t <= 1000; t++ {
		for i := 0; i < n; i++ {
			if t >= b[i] && (t-b[i])%a[i] == 0 {
				if state[i] == 1 {
					state[i] = 0
					cur--
				} else {
					state[i] = 1
					cur++
				}
			}
		}
		if cur > maxOn {
			maxOn = cur
		}
	}

	fmt.Println(maxOn)
}
