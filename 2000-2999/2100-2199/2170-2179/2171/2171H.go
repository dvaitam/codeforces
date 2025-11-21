package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var scores [3]int
	for i := range scores {
		if _, err := fmt.Fscan(in, &scores[i]); err != nil {
			return
		}
	}

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
		fmt.Println("check again")
		return
	}

	s := scores[:]
	sort.Ints(s)
	fmt.Printf("final %d\n", s[1])
}

