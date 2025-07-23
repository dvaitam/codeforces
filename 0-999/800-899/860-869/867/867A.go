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
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	sf := 0
	fs := 0
	for i := 0; i+1 < n; i++ {
		if s[i] == 'S' && s[i+1] == 'F' {
			sf++
		} else if s[i] == 'F' && s[i+1] == 'S' {
			fs++
		}
	}
	if sf > fs {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
