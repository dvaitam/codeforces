package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		var s string
		fmt.Fscan(reader, &s)

		lower := make([]int, 26)
		upper := make([]int, 26)
		for _, ch := range s {
			if ch >= 'a' && ch <= 'z' {
				lower[ch-'a']++
			} else if ch >= 'A' && ch <= 'Z' {
				upper[ch-'A']++
			}
		}

		ans := 0
		for i := 0; i < 26; i++ {
			pairs := min(lower[i], upper[i])
			ans += pairs
			diff := abs(lower[i] - upper[i])
			extra := min(k, diff/2)
			ans += extra
			k -= extra
		}
		fmt.Fprintln(writer, ans)
	}
}
