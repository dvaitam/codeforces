package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var m int
	var k int64
	if _, err := fmt.Fscan(reader, &m, &k); err != nil {
		return
	}
	freq := make([]int64, m+1)
	for i := 0; i <= m; i++ {
		fmt.Fscan(reader, &freq[i])
	}
	ans := 0
	for i := 1; ; i++ {
		var f int64
		if i <= m {
			f = freq[i]
		} else {
			f = 0
		}
		if f <= int64(i)*k {
			ans = i
			break
		}
	}
	fmt.Fprintln(writer, ans)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solve(reader, writer)
	}
}
