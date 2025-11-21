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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		sets := make([][]int, n)
		elemSets := make([][]int, m+1)
		for i := 0; i < n; i++ {
			var l int
			fmt.Fscan(in, &l)
			sets[i] = make([]int, l)
			for j := 0; j < l; j++ {
				fmt.Fscan(in, &sets[i][j])
				elem := sets[i][j]
				elemSets[elem] = append(elemSets[elem], i)
			}
		}

		if canHaveThree(sets, elemSets, m) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

func canHaveThree(sets [][]int, elemSets [][]int, m int) bool {
	n := len(sets)
	freq := make([]int, m+1)
	active := make([]bool, m+1)
	remaining := 0

	for x := 1; x <= m; x++ {
		freq[x] = len(elemSets[x])
		if freq[x] == 0 {
			return false
		}
		active[x] = true
		remaining++
	}

	forced := make([]bool, n)
	queue := make([]int, 0)
	for x := 1; x <= m; x++ {
		if freq[x] == 1 {
			queue = append(queue, x)
		}
	}

	forcedCount := 0
	for head := 0; head < len(queue); head++ {
		x := queue[head]
		if !active[x] {
			continue
		}
		sid := -1
		for _, idx := range elemSets[x] {
			if !forced[idx] {
				sid = idx
				break
			}
		}
		if sid == -1 {
			continue
		}
		forced[sid] = true
		forcedCount++
		for _, y := range sets[sid] {
			if active[y] {
				active[y] = false
				remaining--
			}
		}
	}

	if remaining > 0 {
		return true
	}

	optional := n - forcedCount
	return optional >= 2
}
