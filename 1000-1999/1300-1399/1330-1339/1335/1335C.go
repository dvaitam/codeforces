package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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
		freq := make(map[int]int)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			freq[x]++
		}
		m := len(freq)
		maxF := 0
		for _, v := range freq {
			if v > maxF {
				maxF = v
			}
		}
		a := min(maxF-1, m)
		b := min(maxF, m-1)
		if b > a {
			a = b
		}
		if a < 0 {
			a = 0
		}
		fmt.Fprintln(writer, a)
	}
}
