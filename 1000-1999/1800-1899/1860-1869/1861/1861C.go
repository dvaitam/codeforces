package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemC.txt.
// Given a sequence of operations on an array ('+' push, '-' pop and
// checks '0'/'1' for unsorted/sorted), it determines whether the
// sequence is consistent with some choice of integers appended and
// removed. The algorithm tracks the current size of the array and
// minimal pops required after the last '0' to restore sorted order.
// It also ensures that a '0' check can only occur after at least one
// element has been added since the last '1'.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)

		size := 0  // current size of the array
		need := -1 // threshold size to drop below to become sorted after a '0'
		added := 0 // how many '+' after the last successful '1' remain in stack
		valid := true

		for i := 0; i < len(s) && valid; i++ {
			switch s[i] {
			case '+':
				size++
				added++
			case '-':
				size--
				if added > 0 {
					added--
				}
				if need != -1 && size < need {
					need = -1
				}
			case '1':
				if need != -1 {
					valid = false
					break
				}
				added = 0
			case '0':
				if size < 2 || added == 0 {
					valid = false
					break
				}
				need = size
			}
		}

		if valid {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
