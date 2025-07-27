package main

import (
	"bufio"
	"fmt"
	"os"
)

func isPalindrome(s string) bool {
	n := len(s)
	for i := 0; i < n/2; i++ {
		if s[i] != s[n-1-i] {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &s)
		zeros := 0
		for i := 0; i < n; i++ {
			if s[i] == '0' {
				zeros++
			}
		}
		if isPalindrome(s) {
			if zeros == 1 {
				fmt.Fprintln(out, "BOB")
			} else if zeros%2 == 1 {
				fmt.Fprintln(out, "ALICE")
			} else {
				fmt.Fprintln(out, "BOB")
			}
		} else {
			diff := 0
			for i := 0; i < n/2; i++ {
				if s[i] != s[n-1-i] {
					diff++
				}
			}
			if diff == 1 && zeros == 2 && n%2 == 1 && s[n/2] == '0' {
				fmt.Fprintln(out, "DRAW")
			} else {
				fmt.Fprintln(out, "ALICE")
			}
		}
	}
}
