package main

import (
	"bufio"
	"fmt"
	"os"
)

const N = 1000000007

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	fmt.Fscan(reader, &n)
	v := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &v[i])
	}
	k := n / 2
	fmt.Fprintln(writer, k)
	for i := 0; i+1 < n; i += 2 {
		a, b := v[i], v[i+1]
		mn := a
		if b < a {
			mn = b
		}
		// output: positions are 1-based
		fmt.Fprintln(writer, i+1, i+2, mn, N)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solve(reader, writer)
	}
}
