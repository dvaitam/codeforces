package main

import (
	"bufio"
	"fmt"
	"os"
)

func countFactors(x int) int {
	cnt := 0
	d := 2
	for d*d <= x {
		for x%d == 0 {
			cnt++
			x /= d
		}
		d++
	}
	if x > 1 {
		cnt++
	}
	return cnt
}

func solve(l, r int) int {
	best := 0
	for x := l; x <= r; x++ {
		if c := countFactors(x); c > best {
			best = c
		}
	}
	return best
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		fmt.Fprintln(writer, solve(l, r))
	}
}
