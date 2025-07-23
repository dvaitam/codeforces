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
		var n int64
		fmt.Fscan(reader, &n)
		total := n * (n + 1) / 2
		var sumPowers int64
		for p := int64(1); p <= n; p <<= 1 {
			sumPowers += p
		}
		result := total - 2*sumPowers
		fmt.Fprintln(writer, result)
	}
}
