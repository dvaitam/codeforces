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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s, z string
		fmt.Fscan(in, &s)
		fmt.Fscan(in, &z)

		if n%2 == 1 {
			if s == z {
				fmt.Fprintln(out, "Yes")
			} else {
				fmt.Fprintln(out, "No")
			}
			continue
		}

		sZero := true
		zZero := true
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				sZero = false
				break
			}
		}
		for i := 0; i < n; i++ {
			if z[i] == '1' {
				zZero = false
				break
			}
		}

		if (sZero && zZero) || (!sZero && !zZero) {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
