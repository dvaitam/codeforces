package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
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
		var n, m int
		fmt.Fscan(in, &n, &m)
		var s string
		fmt.Fscan(in, &s)
		L := n * m
		pref := make([]int, L+1)
		for i := 0; i < L; i++ {
			pref[i+1] = pref[i]
			if s[i] == '1' {
				pref[i+1]++
			}
		}

		rows := make([]int, L+1)
		colSeen := make([]bool, m)
		columns := make([]int, L+1)
		totalCols := 0

		for i := 1; i <= L; i++ {
			r := (i - 1) % m
			if s[i-1] == '1' && !colSeen[r] {
				colSeen[r] = true
				totalCols++
			}
			columns[i] = totalCols

			if pref[i]-pref[max(0, i-m)] > 0 {
				rows[i] = rows[max(0, i-m)] + 1
			} else {
				rows[i] = rows[max(0, i-m)]
			}
		}

		for i := 1; i <= L; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, rows[i]+columns[i])
		}
		fmt.Fprintln(out)
	}
}
