package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(in, &s)
	isPal := true
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
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
