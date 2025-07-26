package main

import (
	"bufio"
	"fmt"
	"os"
)

func query(out *bufio.Writer, in *bufio.Reader, idx []int) int {
	fmt.Fprintf(out, "? %d", len(idx))
	for _, v := range idx {
		fmt.Fprintf(out, " %d", v)
	}
	fmt.Fprintln(out)
	out.Flush()
	var res int
	fmt.Fscan(in, &res)
	return res
}

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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		candidate := make([]int, n)
		for i := 0; i < n; i++ {
			candidate[i] = i + 1
		}

		for len(candidate) > 1 {
			m := len(candidate) / 2
			subset := candidate[:m]
			expected := 0
			for _, idx := range subset {
				expected += a[idx-1]
			}
			res := query(out, in, subset)
			if res > expected {
				candidate = subset
			} else {
				candidate = candidate[m:]
			}
		}

		fmt.Fprintf(out, "! %d\n", candidate[0])
		out.Flush()
	}
}
