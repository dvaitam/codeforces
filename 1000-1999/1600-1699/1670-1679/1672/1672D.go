package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	fmt.Fscan(reader, &n)
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	i := 0
	cnt := make(map[int]int)
	for j := 0; j < n; j++ {
		if i < n && a[i] == b[j] {
			i++
			continue
		}
		if j > 0 && b[j] == b[j-1] && cnt[b[j]] > 0 {
			cnt[b[j]]--
			continue
		}
		for i < n && a[i] != b[j] {
			cnt[a[i]]++
			i++
		}
		if i == n {
			fmt.Fprintln(writer, "NO")
			return
		}
		i++
	}
	fmt.Fprintln(writer, "YES")
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
