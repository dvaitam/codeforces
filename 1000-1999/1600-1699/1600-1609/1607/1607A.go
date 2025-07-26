package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
		var layout, s string
		fmt.Fscan(reader, &layout)
		fmt.Fscan(reader, &s)

		pos := make([]int, 26)
		for i := 0; i < 26; i++ {
			pos[layout[i]-'a'] = i
		}

		time := 0
		for i := 1; i < len(s); i++ {
			time += abs(pos[s[i]-'a'] - pos[s[i-1]-'a'])
		}
		fmt.Fprintln(writer, time)
	}
}
