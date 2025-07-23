package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	cntZero := 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		if x == 0 {
			cntZero++
		}
	}

	if n == 1 {
		if cntZero == 0 {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	} else {
		if cntZero == 1 {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
