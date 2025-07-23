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

	var l, r, k int64
	if _, err := fmt.Fscan(in, &l, &r, &k); err != nil {
		return
	}

	if k == 1 {
		if l <= 1 && 1 <= r {
			fmt.Fprintln(out, 1)
		} else {
			fmt.Fprintln(out, -1)
		}
		return
	}

	var res []int64
	for x := int64(1); x <= r; {
		if x >= l {
			res = append(res, x)
		}
		if x > r/k {
			break
		}
		x *= k
	}

	if len(res) == 0 {
		fmt.Fprintln(out, -1)
		return
	}

	for i, v := range res {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
