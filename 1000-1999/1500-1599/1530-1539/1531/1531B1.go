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

	w := make([]int64, n)
	h := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &w[i], &h[i])
	}

	var ans int64
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if w[i] == w[j] || w[i] == h[j] || h[i] == w[j] || h[i] == h[j] {
				ans++
			}
		}
	}

	fmt.Fprintln(writer, ans)
}
