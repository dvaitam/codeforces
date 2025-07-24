package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement solution for problem G as described in problemG.txt.
// For now this is just a placeholder that reads the input format
// and outputs 0 to allow building and testing of the repository.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}
	fmt.Fprintln(out, 0)
}
