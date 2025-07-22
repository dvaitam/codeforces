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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	colors := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &colors[i])
	}
	need := make([]int, m)
	total := 0
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &need[i])
		total += need[i]
	}

	exists := false
	for start := 0; start+total <= n && !exists; start++ {
		cnt := make([]int, m)
		for j := start; j < start+total; j++ {
			c := colors[j] - 1
			cnt[c]++
		}
		match := true
		for i := 0; i < m; i++ {
			if cnt[i] != need[i] {
				match = false
				break
			}
		}
		if match {
			exists = true
		}
	}

	if exists {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
}
