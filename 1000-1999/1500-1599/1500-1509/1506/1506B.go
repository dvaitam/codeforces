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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		var s string
		fmt.Fscan(reader, &s)
		// find positions of stars
		first := -1
		last := -1
		for i := 0; i < n; i++ {
			if s[i] == '*' {
				if first == -1 {
					first = i
				}
				last = i
			}
		}
		if first == last {
			fmt.Fprintln(writer, 1)
			continue
		}
		count := 1
		pos := first
		for pos < last {
			next := pos + k
			if next > last {
				next = last
			}
			for next > pos && s[next] != '*' {
				next--
			}
			if next == pos {
				break
			}
			count++
			pos = next
		}
		fmt.Fprintln(writer, count)
	}
}
