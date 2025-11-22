package main

import (
	"bufio"
	"fmt"
	"os"
)

// For a fixed number of queries q, we can preassign each number a length-q codeword
// over alphabet {L, R, N} describing where it would be put in each query.
// One answer is ignored, so two numbers would be indistinguishable if their
// codewords differ in at most one position. Hence we need codewords with
// Hamming distance at least 2.
//
// Singleton bound for minimum distance 2 over alphabet size 3 gives
//   max codewords M <= 3^{q-1}.
// Achievable via adding one parity digit (last symbol fixes the sum modulo 3)
// to any (q-1)-length ternary word, which ensures distance >= 2.
// Therefore the optimal number of queries is the smallest q with 3^{q-1} >= n.
//
// So f(n) = ceil(log_3 n) + 1.

func minQueries(n int64) int {
	q := 1
	pow := int64(1) // 3^{q-1}
	for pow < n {
		q++
		pow *= 3
	}
	return q
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int64
		fmt.Fscan(in, &n)
		fmt.Fprintln(out, minQueries(n))
	}
}
