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
	var pattern string
	fmt.Fscan(in, &pattern)
	var m int
	fmt.Fscan(in, &m)
	words := make([]string, 0, m)
	for i := 0; i < m; i++ {
		var w string
		fmt.Fscan(in, &w)
		ok := true
		for j := 0; j < n && ok; j++ {
			if pattern[j] != '*' && pattern[j] != w[j] {
				ok = false
			}
		}
		if ok {
			words = append(words, w)
		}
	}
	guessed := make(map[byte]bool)
	for i := 0; i < n; i++ {
		if pattern[i] != '*' {
			guessed[pattern[i]] = true
		}
	}
	ans := 0
	for ch := byte('a'); ch <= 'z'; ch++ {
		if guessed[ch] {
			continue
		}
		good := true
		for _, w := range words {
			found := false
			for i := 0; i < n; i++ {
				if pattern[i] == '*' && w[i] == ch {
					found = true
					break
				}
			}
			if !found {
				good = false
				break
			}
		}
		if good {
			ans++
		}
	}
	fmt.Fprintln(out, ans)
}
