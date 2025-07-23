package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	m := n / 2
	pos := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &pos[i])
	}
	sort.Ints(pos)

	black := 0
	white := 0
	for i := 0; i < m; i++ {
		black += abs(pos[i] - (2*i + 1))
		white += abs(pos[i] - (2 * (i + 1)))
	}
	if black < white {
		fmt.Fprintln(out, black)
	} else {
		fmt.Fprintln(out, white)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
