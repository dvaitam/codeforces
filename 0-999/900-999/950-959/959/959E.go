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

	m := n - 1
	var res int64
	for bit := int64(1); bit <= m; bit <<= 1 {
		res += bit * (m/bit - m/(bit<<1))
	}
	fmt.Fprintln(out, res)
}
