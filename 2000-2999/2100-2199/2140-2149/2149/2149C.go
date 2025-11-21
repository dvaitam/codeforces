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
		var n, k int
		fmt.Fscan(in, &n, &k)

		present := make([]bool, k)
		cntK := 0

		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x < k {
				present[x] = true
			} else if x == k {
				cntK++
			}
		}

		missing := 0
		for i := 0; i < k; i++ {
			if !present[i] {
				missing++
			}
		}

		if cntK < missing {
			fmt.Fprintln(out, missing)
		} else {
			fmt.Fprintln(out, cntK)
		}
	}
}
