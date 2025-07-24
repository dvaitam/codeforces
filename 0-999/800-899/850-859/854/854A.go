package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	for a := n / 2; a >= 1; a-- {
		b := n - a
		if a < b && gcd(a, b) == 1 {
			fmt.Println(a, b)
			return
		}
	}
}
