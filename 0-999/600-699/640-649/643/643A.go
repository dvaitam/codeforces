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
	colors := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &colors[i])
	}

	ans := make([]int, n+1)
	freq := make([]int, n+1)
	for l := 0; l < n; l++ {
		for i := 1; i <= n; i++ {
			freq[i] = 0
		}
		maxColor := 0
		maxCount := 0
		for r := l; r < n; r++ {
			c := colors[r]
			freq[c]++
			if freq[c] > maxCount || (freq[c] == maxCount && c < maxColor) {
				maxCount = freq[c]
				maxColor = c
			}
			ans[maxColor]++
		}
	}

	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, ans[i])
	}
	fmt.Fprintln(writer)
}
