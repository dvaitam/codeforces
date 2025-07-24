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
	var x int
	if _, err := fmt.Fscan(in, &n, &x); err != nil {
		return
	}
	sum := 0
	for i := 0; i < n; i++ {
		var ai int
		fmt.Fscan(in, &ai)
		sum += ai
	}

	if sum+n-1 == x {
		fmt.Fprintln(out, "YES")
	} else {
		fmt.Fprintln(out, "NO")
	}
}
