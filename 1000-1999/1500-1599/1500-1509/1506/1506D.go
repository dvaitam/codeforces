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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		freq := make(map[int]int)
		maxCnt := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			freq[x]++
			if freq[x] > maxCnt {
				maxCnt = freq[x]
			}
		}
		if 2*maxCnt <= n {
			fmt.Fprintln(writer, n%2)
		} else {
			fmt.Fprintln(writer, 2*maxCnt-n)
		}
	}
}
