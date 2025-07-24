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
		var x int64
		fmt.Fscan(reader, &n, &x)
		var s string
		fmt.Fscan(reader, &s)
		prefix := int64(0)
		prefixes := make([]int64, n)
		for i := 0; i < n; i++ {
			if s[i] == '0' {
				prefix++
			} else {
				prefix--
			}
			prefixes[i] = prefix
		}
		total := prefix
		if total == 0 {
			infinite := x == 0
			for i := 0; i < n && !infinite; i++ {
				if prefixes[i] == x {
					infinite = true
				}
			}
			if infinite {
				fmt.Fprintln(writer, -1)
			} else {
				fmt.Fprintln(writer, 0)
			}
			continue
		}
		ans := 0
		if x%total == 0 && x/total >= 0 {
			ans++
		}
		for i := 0; i < n; i++ {
			diff := x - prefixes[i]
			if diff%total == 0 && diff/total >= 0 {
				ans++
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
