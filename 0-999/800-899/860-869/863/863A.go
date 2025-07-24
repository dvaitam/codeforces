package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	// remove trailing zeros
	for len(s) > 0 && s[len(s)-1] == '0' {
		s = s[:len(s)-1]
	}
	// check palindrome
	isPal := true
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-1-i] {
			isPal = false
			break
		}
	}
	if isPal {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
