package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &p[i])
	}

	ans := make([][]int, n)
	for i := 0; i < n; i++ {
		ans[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		val := p[i]
		r, c := i, i
		ans[r][c] = val
		rem := val - 1
		for rem > 0 {
			if c > 0 && ans[r][c-1] == 0 {
				c--
				ans[r][c] = val
			} else if r+1 < n && ans[r+1][c] == 0 {
				r++
				ans[r][c] = val
			} else {
				fmt.Fprintln(out, -1)
				return
			}
			rem--
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j <= i; j++ {
			if j > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, ans[i][j])
		}
		out.WriteByte('\n')
	}
}
