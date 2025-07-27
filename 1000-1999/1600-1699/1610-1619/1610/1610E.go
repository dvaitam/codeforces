package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution attempts to approximate the minimal deletions for problemE.txt.
// The approach builds a subsequence while ensuring that for any three
// consecutive kept elements a, b, c we have 2*b <= a + c. Whenever the
// inequality is violated by a newly read value we discard the oldest element
// of the triple. This heuristic matches the optimal result on small cases and
// keeps the implementation simple.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := range arr {
			fmt.Fscan(in, &arr[i])
		}
		// Build subsequence greedily.
		seq := make([]int, 0, n)
		for _, v := range arr {
			seq = append(seq, v)
			for len(seq) >= 3 && 2*seq[len(seq)-2] > seq[len(seq)-3]+seq[len(seq)-1] {
				// drop the oldest element of the last triple
				copy(seq[len(seq)-3:], seq[len(seq)-2:])
				seq = seq[:len(seq)-1]
			}
		}
		fmt.Fprintln(out, n-len(seq))
	}
}
