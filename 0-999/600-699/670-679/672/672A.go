package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	// Build concatenated string until length >= n
	s := make([]byte, 0, n+10)
	for i := 1; len(s) < n; i++ {
		s = append(s, fmt.Sprint(i)...)
	}

	fmt.Printf("%c\n", s[n-1])
}
