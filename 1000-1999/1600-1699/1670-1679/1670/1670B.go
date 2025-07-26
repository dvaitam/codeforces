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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)
		var k int
		fmt.Fscan(reader, &k)
		special := make(map[rune]bool)
		for i := 0; i < k; i++ {
			var c string
			fmt.Fscan(reader, &c)
			special[rune(c[0])] = true
		}
		prev := 0
		ans := 0
		for i, ch := range s {
			if special[ch] {
				gap := (i + 1) - prev
				if gap > ans {
					ans = gap
				}
				prev = i + 1
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
