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

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var n int
		fmt.Fscan(in, &n)
		var a, b string
		fmt.Fscan(in, &a)
		fmt.Fscan(in, &b)
		s := [2]string{a, b}
		row := 0
		ok := true
		for i := 0; i < n; i++ {
			d := s[row][i] - '0'
			if d <= 2 {
				// straight pipe, continue in same row
				continue
			}
			// curved pipe, need to switch rows
			d2 := s[1-row][i] - '0'
			if d2 <= 2 {
				ok = false
				break
			}
			row ^= 1
		}
		if ok && row == 1 {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
