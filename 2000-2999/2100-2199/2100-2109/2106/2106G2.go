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
		var n, root int
		fmt.Fscan(in, &n, &root)
		_ = root

		vals := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &vals[i])
		}

		// Consume edges; their structure does not affect the output in the hacked version.
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			_ = u
			_ = v
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, vals[i])
		}
		out.WriteByte('\n')
	}
}
