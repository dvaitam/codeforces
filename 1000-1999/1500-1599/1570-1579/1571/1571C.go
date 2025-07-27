package main

import (
	"bufio"
	"fmt"
	"os"
)

func commonSuffixLength(a, b string) int {
	i := len(a) - 1
	j := len(b) - 1
	l := 0
	for i >= 0 && j >= 0 {
		if a[i] != b[j] {
			break
		}
		l++
		i--
		j--
	}
	return l
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		minRhyme := 1 << 60
		maxNo := -1
		for i := 0; i < n; i++ {
			var s, t string
			var r int
			fmt.Fscan(in, &s, &t, &r)
			l := commonSuffixLength(s, t)
			if r == 1 {
				if l < minRhyme {
					minRhyme = l
				}
			} else {
				if l > maxNo {
					maxNo = l
				}
			}
		}
		start := maxNo + 1
		if start < 0 {
			start = 0
		}
		if start > minRhyme {
			fmt.Fprintln(out, 0)
		} else {
			m := minRhyme - start + 1
			fmt.Fprint(out, m)
			for k := start; k <= minRhyme; k++ {
				fmt.Fprint(out, " ", k)
			}
			fmt.Fprintln(out)
		}
	}
}
