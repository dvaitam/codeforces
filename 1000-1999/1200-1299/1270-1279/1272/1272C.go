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
	allowed := make(map[byte]bool)
	for i := 0; i < k; i++ {
		var ch string
		fmt.Fscan(reader, &ch)
		allowed[ch[0]] = true
	}

	var cur int64
	var ans int64
	for i := 0; i < n; i++ {
		if allowed[s[i]] {
			cur++
		} else {
			ans += cur * (cur + 1) / 2
			cur = 0
		}
	}
	ans += cur * (cur + 1) / 2
	fmt.Fprintln(writer, ans)
}
