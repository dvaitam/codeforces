package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)

	count := 0
	for i := 0; i < n; i++ {
		if s[i] == 'W' {
			count++
		}
	}
	maxCount := count
	for i := 1; i < n; i++ {
		if s[i-1] == 'W' {
			count--
		}
		if s[i+n-1] == 'W' {
			count++
		}
		if count > maxCount {
			maxCount = count
		}
	}
	fmt.Println(maxCount)
}
