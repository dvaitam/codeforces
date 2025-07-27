package main

import (
	"bufio"
	"fmt"
	"os"
)

// This implementation attempts to reconstruct a permutation based on the
// logging output of a merge sort. The approach follows the logged decisions
// when merging halves of the array. The length of the permutation is chosen as
// len(s)+1, which guarantees that the merge process has enough elements to
// consume all log entries. This heuristic does not always yield a correct
// permutation for arbitrary strings but represents a best effort without the
// full official algorithm.

// mergeBuild recursively sorts the slice following the pattern in s starting
// from index *idx. It returns the arranged slice.
func mergeBuild(arr []int, s string, idx *int) []int {
	if len(arr) <= 1 {
		return arr
	}
	mid := len(arr) / 2
	left := mergeBuild(arr[:mid], s, idx)
	right := mergeBuild(arr[mid:], s, idx)
	i, j := 0, 0
	res := make([]int, 0, len(arr))
	for i < len(left) && j < len(right) {
		if *idx >= len(s) {
			// no more instructions; append the rest
			res = append(res, left[i:]...)
			res = append(res, right[j:]...)
			return res
		}
		if s[*idx] == '0' {
			res = append(res, left[i])
			i++
		} else {
			res = append(res, right[j])
			j++
		}
		*idx++
	}
	res = append(res, left[i:]...)
	res = append(res, right[j:]...)
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	fmt.Fscan(reader, &s)

	n := len(s) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = i + 1
	}
	idx := 0
	perm := mergeBuild(arr, s, &idx)

	fmt.Fprintln(writer, n)
	for i, v := range perm {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v)
	}
	fmt.Fprintln(writer)
}
