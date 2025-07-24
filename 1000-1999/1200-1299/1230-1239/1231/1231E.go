package main

import (
	"bufio"
	"fmt"
	"os"
)

func isSubsequence(s string, sub string) bool {
	i := 0
	for j := 0; j < len(sub); j++ {
		c := sub[j]
		for i < len(s) && s[i] != c {
			i++
		}
		if i == len(s) {
			return false
		}
		i++
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var n int
		fmt.Fscan(reader, &n)
		var s, t string
		fmt.Fscan(reader, &s)
		fmt.Fscan(reader, &t)
		if len(s) != n || len(t) != n {
			fmt.Println(-1)
			continue
		}
		cnt := make([]int, 26)
		for i := 0; i < n; i++ {
			cnt[s[i]-'a']++
			cnt[t[i]-'a']--
		}
		ok := true
		for _, v := range cnt {
			if v != 0 {
				ok = false
				break
			}
		}
		if !ok {
			fmt.Println(-1)
			continue
		}
		best := 0
		for l := 0; l < n; l++ {
			for r := l; r < n; r++ {
				if isSubsequence(s, t[l:r+1]) {
					if r-l+1 > best {
						best = r - l + 1
					}
				}
			}
		}
		fmt.Println(n - best)
	}
}
