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

	var g, c, l int
	if _, err := fmt.Fscan(in, &g, &c, &l); err != nil {
		return
	}

	maxScore, minScore := g, g
	if c > maxScore {
		maxScore = c
	}
	if l > maxScore {
		maxScore = l
	}
	if c < minScore {
		minScore = c
	}
	if l < minScore {
		minScore = l
	}

	if maxScore-minScore >= 10 {
		fmt.Fprintln(out, "check again")
		return
	}

	median := g + c + l - maxScore - minScore
	fmt.Fprintf(out, "final %d\n", median)
}
