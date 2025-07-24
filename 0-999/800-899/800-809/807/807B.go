package main

import (
	"bufio"
	"fmt"
	"os"
)

func canGetTShirt(s, p int) bool {
	cur := s
	for i := 0; i < 25; i++ {
		cur = (cur*96 + 42) % 475
		if cur+26 == p {
			return true
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var p, x, y int
	if _, err := fmt.Fscan(in, &p, &x, &y); err != nil {
		return
	}

	for s := x; s >= y; s -= 50 {
		if canGetTShirt(s, p) {
			fmt.Fprintln(out, 0)
			return
		}
	}

	for i := 1; ; i++ {
		s := x + 50*i
		if canGetTShirt(s, p) {
			ans := (i + 1) / 2
			fmt.Fprintln(out, ans)
			return
		}
	}
}
