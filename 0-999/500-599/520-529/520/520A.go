package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var s string
	// Read input length and the string
	fmt.Fscan(reader, &n)
	fmt.Fscan(reader, &s)
	seen := make([]bool, 26)
	count := 0
	for _, ch := range s {
		switch {
		case ch >= 'a' && ch <= 'z':
			idx := ch - 'a'
			if !seen[idx] {
				seen[idx] = true
				count++
			}
		case ch >= 'A' && ch <= 'Z':
			idx := ch - 'A'
			if !seen[idx] {
				seen[idx] = true
				count++
			}
		}
	}
	if count == 26 {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
