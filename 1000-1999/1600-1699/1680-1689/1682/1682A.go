package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves problemA.txt for contest 1682A.
// Given a palindromic string s, it counts how many indices can be removed
// while keeping the string a palindrome. The answer equals the length of the
// central block of identical characters because removing a character outside
// this block breaks the palindrome.
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
		if n%2 == 1 {
			mid := n / 2
			ch := s[mid]
			ans := 1
			i := mid - 1
			for i >= 0 && s[i] == ch {
				ans++
				i--
			}
			j := mid + 1
			for j < n && s[j] == ch {
				ans++
				j++
			}
			fmt.Fprintln(out, ans)
		} else {
			left := n/2 - 1
			right := n / 2
			ch := s[left]
			ans := 0
			i := left
			for i >= 0 && s[i] == ch {
				ans++
				i--
			}
			j := right
			for j < n && s[j] == ch {
				ans++
				j++
			}
			fmt.Fprintln(out, ans)
		}
	}
}
