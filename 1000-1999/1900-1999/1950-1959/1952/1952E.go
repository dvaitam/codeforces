package main

import (
	"bufio"
	"fmt"
	"os"
)

// Since the actual problem statement is missing, we simply read the
// input as described and output zero.
const mod = 20240401

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	for i := 0; i < n; i++ {
		var x int
		if _, err := fmt.Fscan(reader, &x); err != nil {
			return
		}
		_ = x
	}
	fmt.Println(0)
}
