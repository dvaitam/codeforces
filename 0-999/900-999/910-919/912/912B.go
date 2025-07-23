package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k uint64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	if k == 1 {
		fmt.Fprintln(writer, n)
		return
	}
	// For k >= 2, maximum XOR-sum is all ones up to the highest bit of n.
	length := bits.Len64(n)
	ans := uint64(1<<uint(length)) - 1
	fmt.Fprintln(writer, ans)
}
