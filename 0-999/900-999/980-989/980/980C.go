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

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	mapping := make([]int, 256)
	for i := range mapping {
		mapping[i] = -1
	}

	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if mapping[x] == -1 {
			start := x
			for start > 0 && x-(start-1)+1 <= k && mapping[start-1] == -1 {
				start--
			}
			if start > 0 && mapping[start-1] != -1 && x-mapping[start-1]+1 <= k {
				start = mapping[start-1]
			}
			for y := start; y <= x; y++ {
				mapping[y] = start
			}
		}
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, mapping[x])
	}
	writer.WriteByte('\n')
}
