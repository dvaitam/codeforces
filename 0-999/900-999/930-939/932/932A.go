package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(reader, &s)
	n := len(s)
	b := make([]byte, n*2)
	for i := 0; i < n; i++ {
		b[i] = s[i]
		b[2*n-1-i] = s[i]
	}
	fmt.Println(string(b))
}
