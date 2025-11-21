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

	var n, m int
	fmt.Fscan(in, &n, &m)

	seq := make([]int64, n)
	if n == 1 {
		seq[0] = 0
	} else {
		for i := 0; i < n; i++ {
			seq[i] = 1 << uint(i)
		}
	}

	for i, val := range seq {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, val)
	}
	fmt.Fprintln(out)
}
