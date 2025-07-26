package main

// This program solves the problem described in problemD.txt for folder 1920.
// It processes two types of operations on an array and answers queries for
// the k-th element after all operations. The array can grow extremely large
// because one of the operations appends several copies of itself. The program
// avoids materializing the array by keeping track of the length after each
// operation and working backwards using binary search for each query.

import (
	"bufio"
	"fmt"
	"os"
)

type Op struct {
	b int   // type of operation (1 or 2)
	x int64 // value or number of copies
}

const INF int64 = 1e18

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(reader, &n, &q)

		ops := make([]Op, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &ops[i].b, &ops[i].x)
		}

		// prefix lengths after each operation (1-indexed, len[0]=0)
		lens := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			if ops[i-1].b == 1 {
				if lens[i-1] < INF {
					lens[i] = lens[i-1] + 1
					if lens[i] > INF {
						lens[i] = INF
					}
				} else {
					lens[i] = INF
				}
			} else {
				if lens[i-1] == 0 {
					lens[i] = 0
				} else if lens[i-1] >= INF/(ops[i-1].x+1) {
					lens[i] = INF
				} else {
					lens[i] = lens[i-1] * (ops[i-1].x + 1)
					if lens[i] > INF {
						lens[i] = INF
					}
				}
			}
		}

		queries := make([]int64, q)
		for i := 0; i < q; i++ {
			fmt.Fscan(reader, &queries[i])
		}

		for i := 0; i < q; i++ {
			k := queries[i]
			idx := n
			for {
				// binary search among operations [1..idx]
				l, r := 1, idx
				pos := idx
				for l <= r {
					m := (l + r) / 2
					if lens[m] >= k {
						pos = m
						r = m - 1
					} else {
						l = m + 1
					}
				}
				op := ops[pos-1]
				if op.b == 1 {
					fmt.Fprint(writer, op.x)
					if i+1 < q {
						writer.WriteByte(' ')
					}
					break
				}
				k = (k-1)%lens[pos-1] + 1
				idx = pos - 1
			}
		}
		fmt.Fprintln(writer)
	}
}
