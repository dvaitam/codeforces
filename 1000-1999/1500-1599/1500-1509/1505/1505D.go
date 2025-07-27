package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	seen := make(map[int]bool)
	for n > 0 {
		d := n % m
		if seen[d] {
			fmt.Println("NO")
			return
		}
		seen[d] = true
		n /= m
	}
	fmt.Println("YES")
}
