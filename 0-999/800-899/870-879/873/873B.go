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
	var s string
	fmt.Fscan(reader, &s)

	diffPos := make(map[int]int)
	diffPos[0] = 0
	diff := 0
	best := 0

	for i := 1; i <= n; i++ {
		if s[i-1] == '1' {
			diff++
		} else {
			diff--
		}
		if pos, ok := diffPos[diff]; ok {
			if i-pos > best {
				best = i - pos
			}
		} else {
			diffPos[diff] = i
		}
	}

	fmt.Fprintln(writer, best)
}
