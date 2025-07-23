package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: Implement solution for problem F.
// Placeholder implementation prints 0.
func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	var mod int64
	if _, err := fmt.Fscan(in, &n, &m, &mod); err != nil {
		return
	}
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
	}
	fmt.Println(0)
}
