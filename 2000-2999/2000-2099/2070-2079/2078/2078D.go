package main

import (
	"bufio"
	"fmt"
	"os"
)

type Gate struct {
	isMul bool
	val   int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		left := make([]Gate, n+1)
		right := make([]Gate, n+1)
		for i := 1; i <= n; i++ {
			var typ string
			var val int64
			fmt.Fscan(in, &typ, &val)
			left[i] = Gate{isMul: typ == "x", val: val}
			fmt.Fscan(in, &typ, &val)
			right[i] = Gate{isMul: typ == "x", val: val}
		}

		vLeft := make([]int64, n+2)
		vRight := make([]int64, n+2)
		vLeft[n+1] = 1
		vRight[n+1] = 1
		for i := n; i >= 1; i-- {
			best := vLeft[i+1]
			if vRight[i+1] > best {
				best = vRight[i+1]
			}
			if left[i].isMul {
				vLeft[i] = vLeft[i+1] + (left[i].val-1)*best
			} else {
				vLeft[i] = vLeft[i+1]
			}
			if right[i].isMul {
				vRight[i] = vRight[i+1] + (right[i].val-1)*best
			} else {
				vRight[i] = vRight[i+1]
			}
		}

		l, r := int64(1), int64(1)
		for i := 1; i <= n; i++ {
			inc := int64(0)
			if left[i].isMul {
				inc += (left[i].val - 1) * l
			} else {
				inc += left[i].val
			}
			if right[i].isMul {
				inc += (right[i].val - 1) * r
			} else {
				inc += right[i].val
			}
			best := vLeft[i+1]
			if vRight[i+1] > best {
				best = vRight[i+1]
			}
			if vLeft[i+1] > vRight[i+1] {
				l += inc
			} else if vLeft[i+1] < vRight[i+1] {
				r += inc
			} else {
				// tie, split evenly
				half := inc / 2
				l += half
				r += inc - half
			}
		}
		fmt.Fprintln(out, l+r)
	}
}
