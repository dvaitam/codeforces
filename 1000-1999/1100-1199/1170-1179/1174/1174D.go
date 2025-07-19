package main

import (
	"bufio"
	"fmt"
	"os"
)

// f computes the value for position i, adjusting for the avoid threshold
func f(i, avoid int) int {
	ans := i & -i
	if ans >= avoid {
		ans <<= 1
	}
	return ans
}

func main() {
	var n, x int
	if _, err := fmt.Scan(&n, &x); err != nil {
		return
	}
	total := 1 << n
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	if x >= total {
		cnt := total - 1
		fmt.Fprintln(w, cnt)
		for i := 1; i < total; i++ {
			if i > 1 {
				w.WriteByte(' ')
			}
			fmt.Fprint(w, i&-i)
		}
		return
	}

	l := (1 << (n - 1)) - 1
	fmt.Fprintln(w, l)
	if l == 0 {
		return
	}
	avoid := x & -x
	for i := 1; i <= l; i++ {
		if i > 1 {
			w.WriteByte(' ')
		}
		fmt.Fprint(w, f(i, avoid))
	}
	fmt.Fprintln(w)
}
