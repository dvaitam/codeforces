package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int64
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	var k int64 = 1
	for n > k {
		n -= k
		k++
	}
	fmt.Fprintln(out, n)
}
