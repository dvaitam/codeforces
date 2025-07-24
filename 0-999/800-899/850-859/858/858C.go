package main

import (
	"bufio"
	"fmt"
	"os"
)

func isVowel(c byte) bool {
	switch c {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}

	cur := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if len(cur) >= 2 && !isVowel(c) && !isVowel(cur[len(cur)-1]) && !isVowel(cur[len(cur)-2]) && !(c == cur[len(cur)-1] && cur[len(cur)-1] == cur[len(cur)-2]) {
			out.Write(cur)
			out.WriteByte(' ')
			cur = cur[:0]
		}
		cur = append(cur, c)
	}
	out.Write(cur)
	out.WriteByte('\n')
}
