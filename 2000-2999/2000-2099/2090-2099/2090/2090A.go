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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x, y, a int64
		fmt.Fscan(in, &x, &y, &a)
		target := a + 1
		if x >= target {
			fmt.Fprintln(out, "NO")
			continue
		}
		sumPair := x + y
		processedPairs := (target - 1) / sumPair
		depth := processedPairs * sumPair
		if depth+x >= target {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
		}
	}
}
