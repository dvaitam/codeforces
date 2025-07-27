package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

// This is a placeholder solution for problem I of contest 1553.
// The original statement is partially missing from the repository so
// we implement a minimal program that reads the input format and
// prints 0 as a default value. This allows the build to succeed.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	// TODO: implement actual counting logic
	fmt.Fprintln(out, 0)
}
