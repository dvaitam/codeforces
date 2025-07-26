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

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		freq := make(map[int]int)
		zeros := 0
		dup := false
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if x == 0 {
				zeros++
			}
			freq[x]++
			if freq[x] > 1 {
				dup = true
			}
		}
		if zeros > 0 {
			fmt.Fprintln(writer, n-zeros)
		} else if dup {
			fmt.Fprintln(writer, n)
		} else {
			fmt.Fprintln(writer, n+1)
		}
	}
}
