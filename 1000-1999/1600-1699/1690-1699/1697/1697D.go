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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	ans := make([]byte, n)
	// last positions of discovered characters, sorted by index ascending
	var pos []int
	var chars []byte

	for i := 0; i < n; i++ {
		// query type 2 on prefix [1, i+1]
		fmt.Fprintf(out, "? 2 1 %d\n", i+1)
		out.Flush()
		var cnt int
		if _, err := fmt.Fscan(in, &cnt); err != nil {
			return
		}
		if cnt > len(pos) {
			// new character
			fmt.Fprintf(out, "? 1 %d\n", i+1)
			out.Flush()
			var s string
			if _, err := fmt.Fscan(in, &s); err != nil {
				return
			}
			c := s[0]
			ans[i] = c
			pos = append(pos, i+1)
			chars = append(chars, c)
		} else {
			// existing character: binary search
			l, r := 0, len(pos)-1
			for l < r {
				m := (l + r + 1) / 2
				fmt.Fprintf(out, "? 2 %d %d\n", pos[m], i+1)
				out.Flush()
				var x int
				if _, err := fmt.Fscan(in, &x); err != nil {
					return
				}
				if x == len(pos)-m {
					l = m
				} else {
					r = m - 1
				}
			}
			ans[i] = chars[l]
			pos[l] = i + 1
			// maintain order
			for l+1 < len(pos) && pos[l] > pos[l+1] {
				pos[l], pos[l+1] = pos[l+1], pos[l]
				chars[l], chars[l+1] = chars[l+1], chars[l]
				l++
			}
		}
	}

	fmt.Fprintf(out, "! %s\n", string(ans))
}
