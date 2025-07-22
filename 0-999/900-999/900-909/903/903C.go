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

	fmt.Fprintln(writer, maxCnt)
}
