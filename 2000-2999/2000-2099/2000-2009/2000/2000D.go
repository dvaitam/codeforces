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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		prefix := make([]int64, n+1)
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			prefix[i+1] = prefix[i] + x
		}

		var s string
		fmt.Fscan(in, &s)

		firstL := -1
		for i := 0; i < n; i++ {
			if s[i] == 'L' {
				firstL = i
				break
			}
		}

		lastR := -1
		for i := n - 1; i >= 0; i-- {
			if s[i] == 'R' {
				lastR = i
				break
			}
		}

		if firstL == -1 || lastR == -1 || firstL >= lastR {
			fmt.Fprintln(out, 0)
			continue
		}

		result := prefix[lastR+1] - prefix[firstL]
		fmt.Fprintln(out, result)
	}
}
