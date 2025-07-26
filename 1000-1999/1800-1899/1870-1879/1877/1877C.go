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
		var n, m, k int64
		fmt.Fscan(reader, &n, &m, &k)
		var ans int64
		if k == 1 {
			ans = 1
		} else if k == 2 {
			part1 := m
			if n < m {
				part1 = n
			}
			var part2 int64
			if m >= n {
				part2 = m/n - 1
			}
			if part2 < 0 {
				part2 = 0
			}
			ans = part1 + part2
		} else if k == 3 {
			part1 := m
			if n < m {
				part1 = n
			}
			var part2 int64
			if m >= n {
				part2 = m/n - 1
			}
			if part2 < 0 {
				part2 = 0
			}
			ans = m - (part1 + part2)
		} else {
			ans = 0
		}
		fmt.Fprintln(writer, ans)
	}
}
