package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var a1, a2, k1, k2, n int
	if _, err := fmt.Fscan(in, &a1, &a2, &k1, &k2, &n); err != nil {
		return
	}

	// minimal players sent off
	safeCards := a1*(k1-1) + a2*(k2-1)
	minPlayers := 0
	if n > safeCards {
		minPlayers = n - safeCards
	}

	// maximal players sent off
	maxPlayers := 0
	cards := n
	if k1 > k2 {
		k1, k2 = k2, k1
		a1, a2 = a2, a1
	}
	use := cards / k1
	if use > a1 {
		use = a1
	}
	maxPlayers += use
	cards -= use * k1

	use = cards / k2
	if use > a2 {
		use = a2
	}
	maxPlayers += use

	fmt.Fprintf(out, "%d %d\n", minPlayers, maxPlayers)
}
