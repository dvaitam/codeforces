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
	n := len(s)
	i := 0
	for i < n && s[i] == 'a' {
		i++
	}
	a := i
	for i < n && s[i] == 'b' {
		i++
	}
	b := i - a
	for i < n && s[i] == 'c' {
		i++
	}
	c := i - a - b
	if i != n || a == 0 || b == 0 {
		fmt.Println("NO")
		return
	}
	if c == a || c == b {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
