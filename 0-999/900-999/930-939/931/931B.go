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

	var n, a, b int
	if _, err := fmt.Fscan(in, &n, &a, &b); err != nil {
		return
	}

	rounds := 0
	for a != b {
		a = (a + 1) / 2
		b = (b + 1) / 2
		rounds++
	}

	totalRounds := 0
	for tmp := n; tmp > 1; tmp /= 2 {
		totalRounds++
	}

	if rounds == totalRounds {
		fmt.Fprintln(out, "Final!")
	} else {
		fmt.Fprintln(out, rounds)
	}
}
