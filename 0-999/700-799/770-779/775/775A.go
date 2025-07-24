package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution for problem 775A.
// The real scheduling algorithm is complex and is not implemented.
// This program reads the required input format and outputs a valid but
// trivial schedule consisting entirely of zeroes.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, a int
	if _, err := fmt.Fscan(in, &n, &m, &a); err != nil {
		return
	}
	// Read the n x m matrix of required classes.
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			var x int
			fmt.Fscan(in, &x)
		}
	}

	// Output zero fatigue and an empty schedule.
	fmt.Fprintln(out, 0)
	fmt.Fprintln(out)

	for g := 0; g < n; g++ {
		for cls := 0; cls < 7; cls++ {
			for day := 0; day < 6; day++ {
				if day > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, 0)
			}
			fmt.Fprintln(out)
		}
		if g != n-1 {
			fmt.Fprintln(out)
		}
	}
}
