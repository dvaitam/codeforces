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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		var s string
		fmt.Fscan(in, &s)
		sb := []byte(s)

		prefMax := make([]int, n+1)
		for i := 1; i <= n; i++ {
			prefMax[i] = prefMax[i-1]
			if p[i-1] > prefMax[i] {
				prefMax[i] = p[i-1]
			}
		}

		good := make([]bool, n+1)
		for i := 1; i < n; i++ {
			if prefMax[i] == i {
				good[i] = true
			}
		}

		blocked := make([]bool, n+1)
		countBad := 0
		for i := 1; i < n; i++ {
			blocked[i] = (sb[i-1] == 'L' && sb[i] == 'R')
			if blocked[i] && !good[i] {
				countBad++
			}
		}

		update := func(i int) {
			if i < 1 || i >= n {
				return
			}
			if blocked[i] && !good[i] {
				countBad--
			}
			blocked[i] = (sb[i-1] == 'L' && sb[i] == 'R')
			if blocked[i] && !good[i] {
				countBad++
			}
		}

		for ; q > 0; q-- {
			var pos int
			fmt.Fscan(in, &pos)
			if sb[pos-1] == 'L' {
				sb[pos-1] = 'R'
			} else {
				sb[pos-1] = 'L'
			}
			update(pos - 1)
			update(pos)
			if countBad == 0 {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		}
	}
}
