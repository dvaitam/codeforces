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
		present := make([]bool, 26)
		for i := 0; i < n; i++ {
			idx := int(s[i] - 'a')
			if idx < 0 || idx >= 26 || present[idx] {
				present = nil
				break
			}
			present[idx] = true
		}
		valid := true
		if present == nil {
			valid = false
		} else {
			for i := 0; i < n; i++ {
				if !present[i] {
					valid = false
					break
				}
			}
		}
		if !valid {
			fmt.Fprintln(out, "NO")
			continue
		}
		l, r := 0, n-1
		ch := byte('a' + n - 1)
		for ch >= 'a' {
			if l <= r && s[l] == ch {
				l++
			} else if l <= r && s[r] == ch {
				r--
			} else {
				valid = false
				break
			}
			ch--
		}
		if valid {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
