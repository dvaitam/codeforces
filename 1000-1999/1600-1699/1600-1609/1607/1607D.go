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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		var s string
		fmt.Fscan(in, &s)

		blues := make([]int, 0, n)
		reds := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if s[i] == 'B' {
				blues = append(blues, a[i])
			} else {
				reds = append(reds, a[i])
			}
		}

		sort.Ints(blues)
		sort.Sort(sort.Reverse(sort.IntSlice(reds)))

		pos := 1
		ok := true
		for _, v := range blues {
			if v < pos {
				ok = false
				break
			}
			pos++
		}

		pos = n
		if ok {
			for _, v := range reds {
				if v > pos {
					ok = false
					break
				}
				pos--
			}
		}

		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
