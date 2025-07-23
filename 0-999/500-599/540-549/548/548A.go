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

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return
	}
	n := len(s)
	// total length must be divisible by k
	if k <= 0 || n%k != 0 {
		fmt.Println("NO")
		return
	}
	m := n / k
	for i := 0; i < k; i++ {
		start := i * m
		end := start + m
		if !isPalindrome(s[start:end]) {
			fmt.Println("NO")
			return
		}
	}
	fmt.Println("YES")
}
