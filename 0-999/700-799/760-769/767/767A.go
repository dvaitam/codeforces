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
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arrived := make([]bool, n+1)
	expected := n
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		arrived[x] = true
		first := true
		for expected > 0 && arrived[expected] {
			if !first {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, expected)
			expected--
			first = false
		}
		writer.WriteByte('\n')
	}
}
