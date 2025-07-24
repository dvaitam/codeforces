package main

import (
	"bufio"
	"fmt"
	"os"
)

func feasible(b int, pos []int, k int) bool {
	maxDiff := b + 1
	prev := pos[0]
	used := 1
	i := 1
	last := pos[len(pos)-1]
	for {
		if last-prev <= maxDiff {
			used++
			return used <= k
		}
		j := i
		for j < len(pos) && pos[j]-prev <= maxDiff {
			j++
		}
		if j == i {
			return false
		}
		prev = pos[j-1]
		used++
		if used > k {
			return false
		}
		i = j
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	pos := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if s[i] == '0' {
			pos = append(pos, i+1)
		}
	}

	left, right := 0, n
	for left < right {
		mid := (left + right) / 2
		if feasible(mid, pos, k) {
			right = mid
		} else {
			left = mid + 1
		}
	}
	fmt.Fprintln(writer, left)
}
