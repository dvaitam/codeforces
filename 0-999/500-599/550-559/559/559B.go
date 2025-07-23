package main

import (
	"bufio"
	"fmt"
	"os"
)

func canonical(s string) string {
	if len(s)%2 == 1 {
		return s
	}
	mid := len(s) / 2
	left := canonical(s[:mid])
	right := canonical(s[mid:])
	if left < right {
		return left + right
	}
	if left > right {
		return right + left
	}
	// left == right
	return left + right
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var a, b string
	if _, err := fmt.Fscan(reader, &a); err != nil {
		return
	}
	fmt.Fscan(reader, &b)
	if canonical(a) == canonical(b) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
