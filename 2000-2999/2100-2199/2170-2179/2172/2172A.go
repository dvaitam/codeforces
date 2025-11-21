package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var g, c, l int
	if _, err := fmt.Fscan(in, &g, &c, &l); err != nil {
		return
	}

	scores := []int{g, c, l}
	minScore, maxScore := scores[0], scores[0]
	for _, v := range scores[1:] {
		if v < minScore {
			minScore = v
		}
		if v > maxScore {
			maxScore = v
		}
	}

	if maxScore-minScore >= 10 {
		fmt.Fprint(out, "check again")
		return
	}

	sort.Ints(scores)
	fmt.Fprintf(out, "final %d", scores[1])
}

