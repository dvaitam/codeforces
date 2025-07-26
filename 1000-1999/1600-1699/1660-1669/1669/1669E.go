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
		cnt := make(map[string]int)
		var ans int64
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(reader, &s)
			// Count pairs differing in one position with previously seen strings
			// iterate over both positions and try changing the letter to others
			for pos := 0; pos < 2; pos++ {
				b := []byte(s)
				orig := b[pos]
				for ch := byte('a'); ch <= byte('k'); ch++ {
					if ch == orig {
						continue
					}
					b[pos] = ch
					ans += int64(cnt[string(b)])
				}
				b[pos] = orig
			}
			cnt[s]++
		}
		fmt.Fprintln(writer, ans)
	}
}
