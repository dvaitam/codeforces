package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)

	alphabet := make([]byte, m)
	for i := 0; i < m; i++ {
		alphabet[i] = byte('a' + i)
	}

	unique := make(map[string]struct{})

	for i := 0; i < n; i++ {
		// string after removing position i
		u := s[:i] + s[i+1:]
		for j := 0; j < n; j++ {
			for _, c := range alphabet {
				if j == i && byte(c) == s[i] {
					continue
				}
				t := u[:j] + string(c) + u[j:]
				if t == s {
					continue
				}
				unique[t] = struct{}{}
			}
		}
	}

	fmt.Println(len(unique))
}
