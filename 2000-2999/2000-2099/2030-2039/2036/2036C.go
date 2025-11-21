package main

import (
	"bufio"
	"fmt"
	"os"
)

func is1100(s []byte, idx int) bool {
	return s[idx] == '1' && s[idx+1] == '1' && s[idx+2] == '0' && s[idx+3] == '0'
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var str string
		fmt.Fscan(in, &str)
		s := []byte(str)
		n := len(s)

		count := 0
		for i := 0; i+3 < n; i++ {
			if is1100(s, i) {
				count++
			}
		}

		var q int
		fmt.Fscan(in, &q)
		for ; q > 0; q-- {
			var pos int
			var v int
			fmt.Fscan(in, &pos, &v)
			idx := pos - 1

			for j := idx - 3; j <= idx; j++ {
				if j >= 0 && j+3 < n && is1100(s, j) {
					count--
				}
			}

			if v == 0 {
				s[idx] = '0'
			} else {
				s[idx] = '1'
			}

			for j := idx - 3; j <= idx; j++ {
				if j >= 0 && j+3 < n && is1100(s, j) {
					count++
				}
			}

			if count > 0 {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		}
	}
}

