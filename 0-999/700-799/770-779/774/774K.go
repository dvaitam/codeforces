package main

import (
	"bufio"
	"fmt"
	"os"
)

func isVowel(c byte) bool {
	switch c {
	case 'a', 'e', 'i', 'o', 'u', 'y':
		return true
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var s string
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	fmt.Fscan(in, &s)

	var res []byte
	i := 0
	for i < n {
		ch := s[i]
		j := i + 1
		for j < n && s[j] == ch {
			j++
		}
		runLen := j - i
		if isVowel(ch) {
			if (ch == 'e' || ch == 'o') && runLen == 2 {
				res = append(res, ch, ch)
			} else {
				res = append(res, ch)
			}
		} else {
			for k := 0; k < runLen; k++ {
				res = append(res, ch)
			}
		}
		i = j
	}

	fmt.Println(string(res))
}
