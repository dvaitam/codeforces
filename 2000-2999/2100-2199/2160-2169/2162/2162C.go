package main

import (
	"bufio"
	"fmt"
	"math/bits"
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
		var a, b int64
		fmt.Fscan(in, &a, &b)

		if a == b {
			fmt.Fprintln(out, 0)
			continue
		}

		if msb(b) > msb(a) {
			fmt.Fprintln(out, -1)
			continue
		}

		cur := a
		target := b
		ops := make([]int64, 0, 2)

		s := (int64(1) << (msb(cur) + 1)) - 1
		if cur != s {
			x := cur ^ s
			ops = append(ops, x)
			cur ^= x
		}

		if cur != target {
			x := cur ^ target
			ops = append(ops, x)
			cur ^= x
		}

		fmt.Fprintln(out, len(ops))
		if len(ops) > 0 {
			for i, v := range ops {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, v)
			}
			fmt.Fprintln(out)
		}
	}
}

func msb(x int64) int {
	return bits.Len64(uint64(x)) - 1
}
