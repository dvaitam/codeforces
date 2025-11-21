package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		rooms := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &rooms[i])
		}
		sort.Ints(rooms)

		low := make([]int, n)
		high := make([]int, n)
		copy(low, rooms[:n])
		copy(high, rooms[m-n:])
		for i := 0; i < n/2; i++ {
			high[i], high[n-1-i] = high[n-1-i], high[i]
		}

		schedule := make([][]int, n)
		for i := 0; i < n; i++ {
			schedule[i] = make([]int, 6)
			for k := 0; k < 6; k++ {
				if k%2 == 0 {
					schedule[i][k] = low[i]
				} else {
					schedule[i][k] = high[i]
				}
			}
		}

		for i := 0; i < n; i++ {
			for k := 0; k < 6; k++ {
				if k > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, schedule[i][k])
			}
			fmt.Fprintln(out)
		}
	}
}
