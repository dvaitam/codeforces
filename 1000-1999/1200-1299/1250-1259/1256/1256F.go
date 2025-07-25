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
		var s, t string
		fmt.Fscan(in, &s)
		fmt.Fscan(in, &t)

		cntS := [26]int{}
		cntT := [26]int{}
		for i := 0; i < n; i++ {
			cntS[s[i]-'a']++
			cntT[t[i]-'a']++
		}
		equal := true
		dup := false
		for i := 0; i < 26; i++ {
			if cntS[i] != cntT[i] {
				equal = false
				break
			}
			if cntS[i] >= 2 {
				dup = true
			}
		}
		if !equal {
			fmt.Fprintln(out, "NO")
			continue
		}
		if dup {
			fmt.Fprintln(out, "YES")
			continue
		}
		// All characters are unique. n <= 26.
		pos := make([]int, 26)
		for i := 0; i < n; i++ {
			pos[t[i]-'a'] = i
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = pos[s[i]-'a']
		}
		invParity := 0
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if arr[i] > arr[j] {
					invParity ^= 1
				}
			}
		}
		if invParity == 0 {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
