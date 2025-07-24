package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	votes := make([][]int, m)
	for i := 0; i < m; i++ {
		votes[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(reader, &votes[i][j])
		}
	}

	totals := make([]int, n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			totals[j] += votes[i][j]
		}
	}

	bestK := m + 1
	var bestIdx []int

	for cand := 0; cand < n-1; cand++ {
		diffTotal := totals[n-1] - totals[cand]
		if diffTotal <= 0 {
			bestK = 0
			bestIdx = nil
			break
		}
		type pair struct{ d, idx int }
		arr := make([]pair, 0, m)
		for i := 0; i < m; i++ {
			diff := votes[i][n-1] - votes[i][cand]
			if diff > 0 {
				arr = append(arr, pair{diff, i})
			}
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i].d > arr[j].d })
		sum := 0
		sel := make([]int, 0)
		for _, p := range arr {
			sum += p.d
			sel = append(sel, p.idx)
			if sum >= diffTotal {
				break
			}
		}
		if sum >= diffTotal && len(sel) < bestK {
			bestK = len(sel)
			bestIdx = append([]int(nil), sel...)
		}
	}

	fmt.Fprintln(writer, bestK)
	for i, id := range bestIdx {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, id+1)
	}
	if len(bestIdx) > 0 {
		fmt.Fprintln(writer)
	}
}
