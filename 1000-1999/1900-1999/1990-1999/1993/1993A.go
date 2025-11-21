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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &s)
		count := map[byte]int{'A': 0, 'B': 0, 'C': 0, 'D': 0}
		for i := 0; i < len(s); i++ {
			if _, ok := count[s[i]]; ok {
				count[s[i]]++
			}
		}
		ans := min(count['A'], n) + min(count['B'], n) + min(count['C'], n) + min(count['D'], n)
		fmt.Fprintln(out, ans)
	}
}
