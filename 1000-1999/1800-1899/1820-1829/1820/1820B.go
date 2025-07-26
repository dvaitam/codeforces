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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		n := len(s)
		tstr := s + s
		maxRun := 0
		cur := 0
		for i := 0; i < len(tstr); i++ {
			if tstr[i] == '1' {
				cur++
				if cur > maxRun {
					maxRun = cur
				}
			} else {
				cur = 0
			}
		}
		if maxRun == 2*n {
			fmt.Fprintln(out, n*n)
			continue
		}
		ans := 0
		if maxRun > n {
			maxRun = n
		}
		for w := 1; w <= maxRun; w++ {
			h := (maxRun + 1 - w)
			if h > n {
				h = n
			}
			area := w * h
			if area > ans {
				ans = area
			}
		}
		fmt.Fprintln(out, ans)
	}
}
