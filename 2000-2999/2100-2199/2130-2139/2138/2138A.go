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
		var k int
		var x uint64
		fmt.Fscan(in, &k, &x)

		total := uint64(1) << (uint(k) + 1)
		mid := total >> 1
		ops := make([]int, 0)

		for x != mid {
			if x < mid {
				x *= 2
				ops = append(ops, 1)
			} else {
				x = 2*x - total
				ops = append(ops, 2)
			}
		}

		for i, j := 0, len(ops)-1; i < j; i, j = i+1, j-1 {
			ops[i], ops[j] = ops[j], ops[i]
		}

		fmt.Fprintln(out, len(ops))
		if len(ops) > 0 {
			for i, op := range ops {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, op)
			}
			fmt.Fprintln(out)
		}
	}
}
