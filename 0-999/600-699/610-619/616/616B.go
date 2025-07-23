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
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	best := int64(-1)
	for i := 0; i < n; i++ {
		minVal := int64(1<<63 - 1)
		for j := 0; j < m; j++ {
			var x int64
			fmt.Fscan(in, &x)
			if x < minVal {
				minVal = x
			}
		}
		if minVal > best {
			best = minVal
		}
	}

	fmt.Fprintln(out, best)
}
