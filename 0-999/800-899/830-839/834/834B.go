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

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	first := make([]int, 26)
	last := make([]int, 26)
	for i := 0; i < 26; i++ {
		first[i] = -1
		last[i] = -1
	}
	for i, ch := range s {
		idx := int(ch - 'A')
		if first[idx] == -1 {
			first[idx] = i
		}
		last[idx] = i
	}

	open := make([]bool, 26)
	openCount := 0
	for i, ch := range s {
		idx := int(ch - 'A')
		if i == first[idx] {
			if !open[idx] {
				open[idx] = true
				openCount++
			}
		}
		if openCount > k {
			fmt.Fprintln(writer, "YES")
			return
		}
		if i == last[idx] {
			if open[idx] {
				open[idx] = false
				openCount--
			}
		}
	}
	fmt.Fprintln(writer, "NO")
}
