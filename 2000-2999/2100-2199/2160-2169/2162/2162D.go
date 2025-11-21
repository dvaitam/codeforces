package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program implements a strategy for the interactive Codeforces problem
// 2162D "Beautiful Permutation".  The interactor stores an unknown
// permutation p and an unknown interval [l, r].  Inside that interval every
// value of p is increased by one, producing array a.  We may query sums over
// arbitrary ranges either in p (query type 1) or in a (query type 2) and need
// to recover l and r in at most 40 queries.
//
// Observations:
//   - The sum of the entire permutation is fixed: n(n+1)/2.
//   - Therefore one query of type 2 over the whole array immediately reveals
//     the length of the modified segment because the total sum increases by
//     exactly (r-l+1).
//   - For any prefix [1, m], the difference between the answers to query type
//     2 and query type 1 equals the number of indices of the hidden segment
//     that fall inside that prefix.  This difference jumps from 0 to a
//     positive value at the left endpoint l, which allows a binary search.
//
// Once l is known, r follows from the previously computed length.
//
// The solution below performs:
//   - 1 query to get the length of the hidden range, and
//   - up to 2*ceil(log2 n) queries during the binary search for l.
//
// This keeps the total below the allowed 40 queries.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return
		}
		totalPermutation := int64(n) * int64(n+1) / 2
		totalModified := query(reader, writer, 2, 1, n)
		segLen := int(totalModified - totalPermutation)

		left, right := 1, n
		for left < right {
			mid := (left + right) / 2
			sumOriginal := query(reader, writer, 1, 1, mid)
			sumModified := query(reader, writer, 2, 1, mid)
			if sumModified-sumOriginal == 0 {
				left = mid + 1
			} else {
				right = mid
			}
		}
		l := left
		r := l + segLen - 1
		fmt.Fprintf(writer, "! %d %d\n", l, r)
		writer.Flush()
	}
}

func query(reader *bufio.Reader, writer *bufio.Writer, typ, l, r int) int64 {
	fmt.Fprintf(writer, "%d %d %d\n", typ, l, r)
	writer.Flush()
	var res int64
	if _, err := fmt.Fscan(reader, &res); err != nil {
		os.Exit(0)
	}
	if res == -1 {
		os.Exit(0)
	}
	return res
}
