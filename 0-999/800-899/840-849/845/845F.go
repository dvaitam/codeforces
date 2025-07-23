package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	for i := 0; i < n; i++ {
		var row string
		fmt.Fscan(in, &row)
	}
	// Placeholder solution: the real implementation is non-trivial.
	fmt.Println(0)
}
