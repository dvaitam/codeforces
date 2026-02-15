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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		cnt := (r + 1) / 2 - l / 2
		fmt.Fprintln(out, cnt / 2)
	}
}