package main

import (
	"bufio"
	"fmt"
	"os"
)

func isPalindrome(s string) bool {
	i, j := 0, len(s)-1
	for i < j {
		if s[i] != s[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func allSame(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}

func solveCase(s string) (bool, []string) {
	n := len(s)
	if n == 1 {
		return false, nil
	}
	if !isPalindrome(s) {
		return true, []string{s}
	}
	if allSame(s) {
		return false, nil
	}
	if n%2 == 0 {
		for _, j := range []int{2, 3} {
			if j <= n-2 {
				left, right := s[:j], s[j:]
				if !isPalindrome(left) && !isPalindrome(right) {
					return true, []string{left, right}
				}
			}
		}
		for j := 2; j <= n-2 && j <= 6; j++ {
			left, right := s[:j], s[j:]
			if !isPalindrome(left) && !isPalindrome(right) {
				return true, []string{left, right}
			}
		}
		return false, nil
	}
	// odd length palindrome
	if n >= 2 {
		x := s[0]
		y := s[1]
		if x != y {
			alt := true
			for i := 0; i < n; i++ {
				if i%2 == 0 && s[i] != x {
					alt = false
					break
				}
				if i%2 == 1 && s[i] != y {
					alt = false
					break
				}
			}
			if alt {
				return false, nil
			}
		}
	}
	outerChar := s[0]
	outerSame := true
	mid := n / 2
	for i := 0; i < n; i++ {
		if i == mid {
			continue
		}
		if s[i] != outerChar {
			outerSame = false
			break
		}
	}
	if outerSame && s[mid] != outerChar {
		return false, nil
	}
	for _, j := range []int{2, 3} {
		if j <= n-2 {
			left, right := s[:j], s[j:]
			if !isPalindrome(left) && !isPalindrome(right) {
				return true, []string{left, right}
			}
		}
	}
	for j := 2; j <= n-2 && j <= 6; j++ {
		left, right := s[:j], s[j:]
		if !isPalindrome(left) && !isPalindrome(right) {
			return true, []string{left, right}
		}
	}
	return false, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var s string
		fmt.Fscan(reader, &s)
		ok, parts := solveCase(s)
		if !ok {
			fmt.Fprintln(writer, "NO")
			continue
		}
		fmt.Fprintln(writer, "YES")
		fmt.Fprintln(writer, len(parts))
		for i, p := range parts {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, p)
		}
		fmt.Fprintln(writer)
	}
}
