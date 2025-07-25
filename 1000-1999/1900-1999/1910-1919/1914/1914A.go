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
		var n int
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)
		counts := make([]int, 26)
		for i := 0; i < len(s); i++ {
			c := s[i] - 'A'
			if c < 26 {
				counts[c]++
			}
		}
		solved := 0
		for i := 0; i < 26; i++ {
			if counts[i] >= i+1 {
				solved++
			}
		}
		fmt.Fprintln(writer, solved)
	}
}
