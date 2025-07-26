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

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		types := make(map[int]struct{})
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			types[x] = struct{}{}
		}
		m := len(types)
		for k := 1; k <= n; k++ {
			if k > 1 {
				fmt.Fprint(out, " ")
			}
			if m > k {
				fmt.Fprint(out, m)
			} else {
				fmt.Fprint(out, k)
			}
		}
		fmt.Fprintln(out)
	}
}
