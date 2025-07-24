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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(reader, &n, &m, &k)
		maxCount := 0
		countMax := 0
		for i := 0; i < m; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if x > maxCount {
				maxCount = x
				countMax = 1
			} else if x == maxCount {
				countMax++
			}
		}
		// condition derived from scheduling theory
		// minimal length to place all occurrences without conflicts
		if int64(maxCount-1)*int64(k)+int64(countMax) <= int64(n) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
