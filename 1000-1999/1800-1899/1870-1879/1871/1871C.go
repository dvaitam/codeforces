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
	n := len(s)
	ok := true
	for i := 0; i < n/2; i++ {
		if s[i] != s[n-1-i] {
			ok = false
			break
		}
	}
	if ok {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
