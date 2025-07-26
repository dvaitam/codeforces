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
		var n, m int
		fmt.Fscan(reader, &n, &m)
		words := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &words[i])
		}
		minDiff := int(^uint(0) >> 1) // max int
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				diff := 0
				for k := 0; k < m; k++ {
					a := words[i][k]
					b := words[j][k]
					if a > b {
						diff += int(a - b)
					} else {
						diff += int(b - a)
					}
				}
				if diff < minDiff {
					minDiff = diff
				}
			}
		}
		fmt.Fprintln(writer, minDiff)
	}
}
