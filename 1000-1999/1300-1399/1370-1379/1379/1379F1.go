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

	var n, m, q int
	if _, err := fmt.Fscan(in, &n, &m, &q); err != nil {
		return
	}

	// arrays indexed from 1..m
	minRow1 := make([]int, m+2) // min row with forced B in column j
	maxRow0 := make([]int, m+2) // max row with forced A in column j
	for j := 1; j <= m; j++ {
		minRow1[j] = n + 1
		maxRow0[j] = 0
	}

	type key struct{ r, c int }
	removed := make(map[key]int)
	failed := false

	prefix := make([]int, m+2)
	suffix := make([]int, m+2)

	for ; q > 0; q-- {
		var i, j int
		fmt.Fscan(in, &i, &j)
		if failed {
			fmt.Fprintln(out, "NO")
			continue
		}
		r := (i + 1) / 2
		c := (j + 1) / 2
		k := key{r, c}
		if i%2 == 1 { // remove A -> force B
			removed[k] |= 1
			if removed[k] == 3 {
				failed = true
			}
			if r < minRow1[c] {
				minRow1[c] = r
			}
		} else {
			removed[k] |= 2
			if removed[k] == 3 {
				failed = true
			}
			if r > maxRow0[c] {
				maxRow0[c] = r
			}
		}
		if failed {
			fmt.Fprintln(out, "NO")
			continue
		}
		// recompute prefix minima and suffix maxima naively
		cur := n + 1
		for jj := 1; jj <= m; jj++ {
			if minRow1[jj] < cur {
				cur = minRow1[jj]
			}
			prefix[jj] = cur
		}
		cur = 0
		for jj := m; jj >= 1; jj-- {
			if maxRow0[jj] > cur {
				cur = maxRow0[jj]
			}
			suffix[jj] = cur
		}
		ok := true
		for jj := 1; jj <= m; jj++ {
			if prefix[jj] <= suffix[jj] {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
