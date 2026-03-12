package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(s, t string) string {
	// Maintain queues of positions for each character
	queues := [26][]int{}
	for i := 0; i < len(s); i++ {
		c := s[i] - 'a'
		queues[c] = append(queues[c], i)
	}
	// Pointers into each queue
	ptr := [26]int{}

	for i := 0; i < len(t); i++ {
		c := t[i] - 'a'
		// Find the next available occurrence of character c
		for ptr[c] < len(queues[c]) && queues[c][ptr[c]] == -1 {
			ptr[c]++
		}
		if ptr[c] >= len(queues[c]) {
			return "NO"
		}
		pos := queues[c][ptr[c]]
		ptr[c]++
		// Remove all characters smaller than c that appear before pos
		for ch := 0; ch < int(c); ch++ {
			for ptr[ch] < len(queues[ch]) && queues[ch][ptr[ch]] < pos {
				ptr[ch]++
			}
		}
	}
	return "YES"
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		var s, t string
		fmt.Fscan(reader, &s)
		fmt.Fscan(reader, &t)
		fmt.Fprintln(writer, solve(s, t))
	}
}
