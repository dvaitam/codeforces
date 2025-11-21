package main

import (
	"bufio"
	"fmt"
	"os"
)

func matches(s []byte, pos int) bool {
	return s[pos] == '1' && s[pos+1] == '1' && s[pos+2] == '0' && s[pos+3] == '0'
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		n := len(s)
		bytes := []byte(s)
		count := 0
		for i := 0; i+3 < n; i++ {
			if matches(bytes, i) {
				count++
			}
		}

		var q int
		fmt.Fscan(in, &q)
		for ; q > 0; q-- {
			var idx, v int
			fmt.Fscan(in, &idx, &v)
			idx--
			target := byte('0' + v)
			if bytes[idx] != target {
				for pos := idx - 3; pos <= idx; pos++ {
					if pos >= 0 && pos+3 < n && matches(bytes, pos) {
						count--
					}
				}
				bytes[idx] = target
				for pos := idx - 3; pos <= idx; pos++ {
					if pos >= 0 && pos+3 < n && matches(bytes, pos) {
						count++
					}
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
