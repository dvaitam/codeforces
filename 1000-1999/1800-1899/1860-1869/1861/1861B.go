package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for problemB.txt from contest 1861.
// Each operation can only remove runs between equal characters.
// Any string can therefore be reduced to a form 0...01...1 with
// some boundary index j marking the last zero. For a fixed string,
// index j is attainable iff every zero after j has a '1' somewhere
// between j and that zero. We precompute for each position the
// smallest index of a left 1 for all later zeros and collect all
// reachable j. Two strings can be made equal iff they share some
// attainable boundary.

func reachableBoundaries(s string) []int {
	n := len(s)
	leftOne := make([]int, n)
	last := 0
	for i := 0; i < n; i++ {
		if s[i] == '1' {
			last = i + 1
		}
		leftOne[i] = last
	}
	const inf = int(1e9)
	minLeft := make([]int, n+1)
	for i := range minLeft {
		minLeft[i] = inf
	}
	for i := n - 1; i >= 0; i-- {
		minLeft[i] = minLeft[i+1]
		if s[i] == '0' {
			if leftOne[i] < minLeft[i] {
				minLeft[i] = leftOne[i]
			}
		}
	}
	res := make([]int, 0)
	for j := 0; j < n-1; j++ {
		if s[j] == '0' && minLeft[j+1] > j+1 {
			res = append(res, j+1)
		}
	}
	return res
}

func canEqual(a, b string) bool {
	ba := reachableBoundaries(a)
	bb := reachableBoundaries(b)
	i, j := 0, 0
	for i < len(ba) && j < len(bb) {
		if ba[i] == bb[j] {
			return true
		}
		if ba[i] < bb[j] {
			i++
		} else {
			j++
		}
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t++ {
		var a, b string
		fmt.Fscan(reader, &a)
		fmt.Fscan(reader, &b)
		if canEqual(a, b) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
