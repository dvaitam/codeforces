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
		var n int
		fmt.Fscan(reader, &n)
		prices := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &prices[i])
		}
		minSuf := prices[n-1]
		bad := 0
		for i := n - 2; i >= 0; i-- {
			if prices[i] > minSuf {
				bad++
			} else {
				minSuf = prices[i]
			}
		}
		fmt.Fprintln(writer, bad)
	}
}
