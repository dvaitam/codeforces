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

	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}
	holes := make([]bool, n+1)
	for i := 0; i < m; i++ {
		var pos int
		fmt.Fscan(in, &pos)
		holes[pos] = true
	}

	position := 1
	if holes[position] {
		fmt.Fprintln(out, position)
		return
	}
	for i := 0; i < k; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		if position == u {
			position = v
		} else if position == v {
			position = u
		}
		if holes[position] {
			fmt.Fprintln(out, position)
			return
		}
	}
	fmt.Fprintln(out, position)
}
