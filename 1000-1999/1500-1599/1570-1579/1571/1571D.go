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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	firstCount := make([]int, n+1)
	lastCount := make([]int, n+1)
	pair := make([][]int, n+1)
	for i := range pair {
		pair[i] = make([]int, n+1)
	}

	var f1, l1 int
	for i := 1; i <= m; i++ {
		var f, l int
		fmt.Fscan(in, &f, &l)
		if i == 1 {
			f1 = f
			l1 = l
		}
		firstCount[f]++
		lastCount[l]++
		pair[f][l]++
	}

	worst := 1
	for x := 1; x <= n; x++ {
		for y := 1; y <= n; y++ {
			if x == y {
				continue
			}
			c2 := pair[x][y]
			c1 := firstCount[x] + lastCount[y] - 2*c2
			var rank int
			if f1 == x && l1 == y {
				rank = 1
			} else if f1 == x || l1 == y {
				rank = c2 + 1
			} else {
				rank = c2 + c1 + 1
			}
			if rank > worst {
				worst = rank
			}
		}
	}

	fmt.Fprintln(out, worst)
}
