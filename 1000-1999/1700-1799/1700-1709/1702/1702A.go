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
		var m int64
		fmt.Fscan(in, &m)
		pow := int64(1)
		for x := m; x >= 10; x /= 10 {
			pow *= 10
		}
		fmt.Fprintln(out, m-pow)
	}
}
