package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &words[i])
	}
	ans := 0
	for ch := byte('a'); ch <= 'e'; ch++ {
		diffs := make([]int, n)
		for i, w := range words {
			cnt := 0
			for j := 0; j < len(w); j++ {
				if w[j] == ch {
					cnt++
				}
			}
			diffs[i] = 2*cnt - len(w)
		}
		sort.Slice(diffs, func(i, j int) bool { return diffs[i] > diffs[j] })
		sum := 0
		cur := 0
		for _, d := range diffs {
			sum += d
			if sum > 0 {
				cur++
			} else {
				break
			}
		}
		if cur > ans {
			ans = cur
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
