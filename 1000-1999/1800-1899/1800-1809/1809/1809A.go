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
		var s string
		fmt.Fscan(reader, &s)
		count := make(map[rune]int)
		for _, ch := range s {
			count[ch]++
		}
		if len(count) == 1 {
			fmt.Fprintln(writer, -1)
			continue
		}
		ans := 4
		for _, c := range count {
			if c == 3 {
				ans = 6
				break
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
