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
		if n%2 == 0 {
			fmt.Fprintln(out, -1)
			continue
		}
		ops := make([]int, 0)
		for n > 1 {
			if n%4 == 1 {
				ops = append(ops, 1)
				n = (n + 1) / 2
			} else {
				ops = append(ops, 2)
				n = (n - 1) / 2
			}
		}
		// reverse operations to get order from start to finish
		for i, j := 0, len(ops)-1; i < j; i, j = i+1, j-1 {
			ops[i], ops[j] = ops[j], ops[i]
		}
		fmt.Fprintln(out, len(ops))
		for i, v := range ops {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
