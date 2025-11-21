package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		if n == 0 {
			fmt.Fprintln(out)
			continue
		}
		if n == 1 {
			fmt.Fprintln(out, "1")
			continue
		}
		fmt.Fprintln(out, strings.Repeat("0", n-1)+"1")
	}
}
