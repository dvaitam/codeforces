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

	var n int64
	var k int
	fmt.Fscan(reader, &n, &k)

	var bestIdx int
	var bestCount int64
	var maxTransport int64

	for i := 1; i <= k; i++ {
		var a int64
		fmt.Fscan(reader, &a)
		cnt := n / a
		transported := cnt * a
		if transported > maxTransport {
			maxTransport = transported
			bestIdx = i
			bestCount = cnt
		}
	}

	fmt.Fprintf(writer, "%d %d\n", bestIdx, bestCount)
}
