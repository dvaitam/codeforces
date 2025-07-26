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
		var n int
		fmt.Fscan(in, &n)
		var s1, s2 string
		fmt.Fscan(in, &s1)
		fmt.Fscan(in, &s2)
		cols := make([]int, 0)
		cnt0, cnt1 := 0, 0
		for i := 0; i < n; i++ {
			if s1[i] == '*' {
				cols = append(cols, i)
				cnt0++
			}
			if s2[i] == '*' {
				cols = append(cols, i)
				cnt1++
			}
		}
		if len(cols) == 1 {
			fmt.Fprintln(out, 0)
			continue
		}
		sort.Ints(cols)
		median := cols[len(cols)/2]
		hor := 0
		for _, c := range cols {
			if c > median {
				hor += c - median
			} else {
				hor += median - c
			}
		}
		if cnt0 < cnt1 {
			hor += cnt0
		} else {
			hor += cnt1
		}
		fmt.Fprintln(out, hor)
	}
}
