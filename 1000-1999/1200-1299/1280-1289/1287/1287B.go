package main

import (
	"bufio"
	"fmt"
	"os"
)

func third(a, b byte) byte {
	if a == b {
		return a
	}
	if a != 'S' && b != 'S' {
		return 'S'
	}
	if a != 'E' && b != 'E' {
		return 'E'
	}
	return 'T'
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	fmt.Fscan(reader, &n, &k)
	cards := make([]string, n)
	index := make(map[string]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &cards[i])
		index[cards[i]] = i
	}

	var count int64
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			b := make([]byte, k)
			for t := 0; t < k; t++ {
				b[t] = third(cards[i][t], cards[j][t])
			}
			if idx, ok := index[string(b)]; ok && idx > j {
				count++
			}
		}
	}

	fmt.Fprintln(writer, count)
}
