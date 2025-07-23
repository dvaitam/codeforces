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
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	res := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if len(res) > 0 && isVowel(c) && isVowel(res[len(res)-1]) {
			continue
		}
		res = append(res, c)
	}
	fmt.Fprintln(writer, string(res))
}
