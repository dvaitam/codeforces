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
		var n, x, p int
		fmt.Fscan(in, &n, &x, &p)
		target := (n - x) % n
		sum := 0
		m := p
		if m > 2*n {
			m = 2 * n
		}
		found := false
		for i := 1; i <= m; i++ {
			sum = (sum + i) % n
			if sum == target {
				found = true
				break
			}
		}
		if found {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
