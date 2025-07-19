package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	f := make([]int, n)
	counts := make(map[int]int)
	cnt, top, loc := 0, 0, 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		f[i] = x
		counts[x]++
		if counts[x] > cnt {
			cnt = counts[x]
			top = x
			loc = i
		}
	}
	// number of operations
	ops := n - cnt
	fmt.Fprintln(writer, ops)
	if ops == 0 {
		return
	}
	// propagate to the right
	for i := loc; i < n; i++ {
		if f[i] < top {
			// type 1: set f[i] = f[i-1]
			fmt.Fprintf(writer, "1 %d %d\n", i+1, i)
		} else if f[i] > top {
			// type 2: set f[i] = f[i-1]
			fmt.Fprintf(writer, "2 %d %d\n", i+1, i)
		}
	}
	// propagate to the left
	for i := loc; i >= 0; i-- {
		if f[i] == top {
			continue
		}
		if f[i] < top {
			// type 1: set f[i] = f[i+1]
			fmt.Fprintf(writer, "1 %d %d\n", i+1, i+2)
		} else if f[i] > top {
			// type 2: set f[i] = f[i+1]
			fmt.Fprintf(writer, "2 %d %d\n", i+1, i+2)
		}
	}
}
