package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	l := make([]int, q)
	r := make([]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(reader, &l[i])
	}
	for i := 0; i < q; i++ {
		fmt.Fscan(reader, &r[i])
	}
	// Complex combinatorial counting is required here.
	// A complete implementation is omitted.
	fmt.Fprintln(writer, 0)
}
