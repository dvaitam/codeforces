package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var a, b int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &a, &b); err != nil {
		return
	}
	// Placeholder solution: simply output number of games in first round modulo mod
	games := n * (n - 1) / 2
	fmt.Println(games % mod)
}
