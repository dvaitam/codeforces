package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for problemB.txt (Ten Words of Wisdom).
// For each test case, among responses with length at most 10 words,
// we choose the one with the highest quality and output its index.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		bestIdx := 1
		bestQual := -1
		for i := 1; i <= n; i++ {
			var a, b int
			fmt.Fscan(reader, &a, &b)
			if a <= 10 && b > bestQual {
				bestQual = b
				bestIdx = i
			}
		}
		fmt.Fprintln(writer, bestIdx)
	}
}
