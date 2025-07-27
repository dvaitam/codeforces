package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n int
	fmt.Fscan(reader, &n)
	// For n=1, trivial permutation
	if n == 1 {
		fmt.Fprintln(writer, "! 1")
		return
	}
	// ans holds the discovered permutation, 1-indexed
	ans := make([]int, n+1)
	// idx is current candidate index for maximum element
	idx := 1
	for i := 2; i <= n; i++ {
		// query p[idx] mod p[i]
		fmt.Fprintf(writer, "? %d %d\n", idx, i)
		writer.Flush()
		var r1 int
		fmt.Fscan(reader, &r1)
		// query p[i] mod p[idx]
		fmt.Fprintf(writer, "? %d %d\n", i, idx)
		writer.Flush()
		var r2 int
		fmt.Fscan(reader, &r2)
		// the smaller element is loser; its value equals mod result
		if r1 > r2 {
			// p[idx] < p[i], so idx is loser
			ans[idx] = r1
			idx = i
		} else {
			// p[i] < p[idx], so i is loser
			ans[i] = r2
		}
	}
	// remaining index holds the maximum element n
	ans[idx] = n
	// output the full permutation
	fmt.Fprint(writer, "! ")
	for i := 1; i <= n; i++ {
		fmt.Fprintf(writer, "%d", ans[i])
		if i < n {
			fmt.Fprint(writer, " ")
		}
	}
	fmt.Fprintln(writer)
}
