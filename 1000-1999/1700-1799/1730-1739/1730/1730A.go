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
		solve(reader, writer)
	}
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n, c int
	fmt.Fscan(reader, &n, &c)
	counts := make(map[int]int)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		counts[x]++
	}
	ans := 0
	for _, cnt := range counts {
		if cnt < c {
			ans += cnt
		} else {
			ans += c
		}
	}
	fmt.Fprintln(writer, ans)
}
