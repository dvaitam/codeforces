package main

import (
	"bufio"
	"fmt"
	"os"
)

func solveCase(s string) string {
	maxCh := byte('a')
	count := 0
	for i := 0; i < len(s); i++ {
		if s[i] > maxCh {
			maxCh = s[i]
			count = 1
		} else if s[i] == maxCh {
			count++
		}
	}
	ans := make([]byte, count)
	for i := 0; i < count; i++ {
		ans[i] = maxCh
	}
	return string(ans)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &s)
		if len(s) != n {
			s = s[:n]
		}
		fmt.Fprint(out, solveCase(s))
		if t > 1 {
			fmt.Fprint(out, "\n")
		}
	}
}
