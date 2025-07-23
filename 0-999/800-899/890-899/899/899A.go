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
	fmt.Fscan(in, &n)
	ones, twos := 0, 0
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(in, &v)
		if v == 1 {
			ones++
		} else {
			twos++
		}
	}

	teams := min(ones, twos)
	ones -= teams
	teams += ones / 3

	fmt.Fprintln(out, teams)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
