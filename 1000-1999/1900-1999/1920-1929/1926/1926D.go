package main

import (
	"bufio"
	"fmt"
	"os"
)

const mask = (1 << 31) - 1

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		freq := make(map[int]int)
		pairs := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			c := mask ^ x
			if freq[c] > 0 {
				freq[c]--
				pairs++
			} else {
				freq[x]++
			}
		}
		fmt.Fprintln(writer, n-pairs)
	}
}
