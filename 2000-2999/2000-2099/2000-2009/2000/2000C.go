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
		var n int
		fmt.Fscan(in, &n)
		pattern := make([]int, n)
		valToID := make(map[int64]int, n)
		nextID := 0
		for i := 0; i < n; i++ {
			var val int64
			fmt.Fscan(in, &val)
			if id, ok := valToID[val]; ok {
				pattern[i] = id
			} else {
				valToID[val] = nextID
				pattern[i] = nextID
				nextID++
			}
		}

		var m int
		fmt.Fscan(in, &m)
		for i := 0; i < m; i++ {
			var s string
			fmt.Fscan(in, &s)
			if len(s) != n {
				fmt.Fprintln(out, "NO")
				continue
			}
			var charID [26]int
			for j := 0; j < 26; j++ {
				charID[j] = -1
			}
			nextCharID := 0
			ok := true
			for idx := 0; idx < n && ok; idx++ {
				c := s[idx] - 'a'
				if c < 0 || c >= 26 {
					ok = false
					break
				}
				if charID[c] == -1 {
					charID[c] = nextCharID
					nextCharID++
				}
				if charID[c] != pattern[idx] {
					ok = false
				}
			}
			if ok {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		}
	}
}
