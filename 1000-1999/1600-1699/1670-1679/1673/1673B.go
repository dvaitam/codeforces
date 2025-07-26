package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		letters := make(map[byte]bool)
		for i := 0; i < len(s); i++ {
			letters[s[i]] = true
		}
		n := len(s)
		pref := make([][26]int, n+1)
		for i := 1; i <= n; i++ {
			pref[i] = pref[i-1]
			pref[i][s[i-1]-'a']++
		}
		pos := make(map[byte][]int)
		for i := 0; i < n; i++ {
			c := s[i]
			pos[c] = append(pos[c], i+1)
		}
		balanced := true
		for _, arr := range pos {
			for j := 1; j < len(arr) && balanced; j++ {
				l := arr[j-1]
				r := arr[j]
				for ch := range letters {
					idx := ch - 'a'
					if pref[r][idx]-pref[l-1][idx] == 0 {
						balanced = false
						break
					}
				}
			}
			if !balanced {
				break
			}
		}
		if balanced {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
