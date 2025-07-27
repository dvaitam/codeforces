package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func solveOne(reader *bufio.Reader, writer *bufio.Writer) {
	var n, k int
	fmt.Fscan(reader, &n, &k)
	pos := make([]int, 0, n)
	neg := make([]int, 0, n)
	maxAbs := 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if x > 0 {
			pos = append(pos, x)
			if x > maxAbs {
				maxAbs = x
			}
		} else if x < 0 {
			neg = append(neg, -x)
			if -x > maxAbs {
				maxAbs = -x
			}
		} else {
			// x == 0, treat as positive or negative doesn't matter since distance 0
		}
	}

	// sort descending
	sort.Slice(pos, func(i, j int) bool { return pos[i] > pos[j] })
	sort.Slice(neg, func(i, j int) bool { return neg[i] > neg[j] })

	total := 0
	for i := 0; i < len(pos); i += k {
		total += pos[i] * 2
	}
	for i := 0; i < len(neg); i += k {
		total += neg[i] * 2
	}

	total -= maxAbs
	fmt.Fprintln(writer, total)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solveOne(reader, writer)
	}
}
