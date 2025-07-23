package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// It reads a sequence of integers and outputs the length of the
// longest contiguous non-decreasing subsequence.
func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	maxLen := 0
	currLen := 0
	var prev int
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		if i == 0 || x >= prev {
			currLen++
		} else {
			currLen = 1
		}
		if currLen > maxLen {
			maxLen = currLen
		}
		prev = x
	}
	fmt.Println(maxLen)
}
