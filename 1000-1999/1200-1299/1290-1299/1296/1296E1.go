package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)

	colors := make([]byte, n)
	last := [2]byte{'a' - 1, 'a' - 1}
	for i := 0; i < n; i++ {
		c := s[i]
		if c >= last[0] {
			colors[i] = '0'
			last[0] = c
		} else if c >= last[1] {
			colors[i] = '1'
			last[1] = c
		} else {
			fmt.Println("NO")
			return
		}
	}
	fmt.Println("YES")
	fmt.Println(string(colors))
}
