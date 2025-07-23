package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	set := make(map[string]struct{})
	for i := 0; i <= len(s); i++ {
		prefix := s[:i]
		suffix := s[i:]
		for c := 'a'; c <= 'z'; c++ {
			str := prefix + string(byte(c)) + suffix
			set[str] = struct{}{}
		}
	}
	fmt.Println(len(set))
}
