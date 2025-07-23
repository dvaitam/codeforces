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
	wins := make([]int, n)
	for i := 0; i < m; i++ {
		bestCand := 0
		var bestVotes int64 = -1
		for j := 0; j < n; j++ {
			var v int64
			fmt.Fscan(in, &v)
			if v > bestVotes {
				bestVotes = v
				bestCand = j
			}
		}
		wins[bestCand]++
	}
	winner := 0
	for i := 1; i < n; i++ {
		if wins[i] > wins[winner] {
			winner = i
		}
	}
	fmt.Println(winner + 1)
}
