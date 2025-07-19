package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var n int
	if _, err := fmt.Fscan(os.Stdin, &n); err != nil {
		return
	}
	ops := make([]int, 0, 3*n*n/2)

	doThing := func(offset, m int) {
		for i := 0; i < m-1; i++ {
			ops = append(ops, (i+offset)%n)
		}
		for k := 1; m-2*k >= 1; k++ {
			for i := m - k - 2; i >= k-1; i-- {
				ops = append(ops, (i+offset)%n)
			}
			for i := k; i <= m-k-1; i++ {
				ops = append(ops, (i+offset)%n)
			}
		}
	}

	doThing(0, n)
	doThing(n-n/2, n-n%2)
	doThing(0, n)

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	fmt.Fprintln(w, len(ops))
	for i, v := range ops {
		if i > 0 {
			w.WriteByte(' ')
		}
		fmt.Fprint(w, v)
	}
	fmt.Fprintln(w)
}
