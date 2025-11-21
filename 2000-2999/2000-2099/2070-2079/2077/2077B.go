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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var x, y, m uint32
		fmt.Fscan(in, &x, &y, &m)
		ans := uint64(x|m) + uint64(y|m)
		fmt.Fprintln(out, ans)
	}
}
