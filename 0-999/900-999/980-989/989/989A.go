package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	n := len(s)
	for i := 0; i+2 < n; i++ {
		hasA, hasB, hasC := false, false, false
		for j := 0; j < 3; j++ {
			switch s[i+j] {
			case 'A':
				hasA = true
			case 'B':
				hasB = true
			case 'C':
				hasC = true
			}
		}
		if hasA && hasB && hasC {
			fmt.Println("Yes")
			return
		}
	}
	fmt.Println("No")
}
