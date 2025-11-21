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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var s, t string
		fmt.Fscan(in, &s)
		fmt.Fscan(in, &t)
		if len(s) != len(t) {
			fmt.Fprintln(out, -1)
			continue
		}
		n := len(s)
		sZeros, tZeros := 0, 0
		for i := 0; i < n; i++ {
			if s[i] == '0' {
				sZeros++
			}
			if t[i] == '0' {
				tZeros++
			}
		}
		if sZeros != tZeros {
			fmt.Fprintln(out, -1)
			continue
		}
		diffZeros := 0
		for i := 0; i < n; i++ {
			if s[i] == '0' && t[i] == '1' {
				diffZeros++
			}
		}
		fmt.Fprintln(out, diffZeros)
	}
}
