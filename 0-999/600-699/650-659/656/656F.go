package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	if len(s) != 7 || s[0] != 'A' {
		return
	}
	prefix := int(s[1]-'0')*10 + int(s[2]-'0')
	hasZero := false
	for i := 3; i < 7; i++ {
		if s[i] == '0' {
			hasZero = true
			break
		}
	}
	if hasZero {
		prefix--
	}
	fmt.Println(prefix)
}
