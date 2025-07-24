package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)
	b := []byte(s)
	for i := 0; i < m; i++ {
		var l, r int
		var c1, c2 string
		fmt.Fscan(reader, &l, &r, &c1, &c2)
		for j := l - 1; j < r; j++ {
			if string(b[j]) == c1 {
				b[j] = c2[0]
			}
		}
	}
	fmt.Println(string(b))
}
