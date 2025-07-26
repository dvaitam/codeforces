package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func computeScorePenalty(times []int, h int) (int, int64) {
	sort.Ints(times)
	var solved int
	var sum int64
	var penalty int64
	for _, t := range times {
		if sum+int64(t) > int64(h) {
			break
		}
		sum += int64(t)
		solved++
		penalty += sum
	}
	return solved, penalty
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, m, h int
		fmt.Fscan(reader, &n, &m, &h)
		participants := make([][]int, n)
		for i := 0; i < n; i++ {
			participants[i] = make([]int, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(reader, &participants[i][j])
			}
		}
		rSolved, rPenalty := computeScorePenalty(participants[0], h)
		rank := 1
		for i := 1; i < n; i++ {
			s, p := computeScorePenalty(participants[i], h)
			if s > rSolved || (s == rSolved && p < rPenalty) {
				rank++
			}
		}
		fmt.Fprintln(writer, rank)
	}
}
