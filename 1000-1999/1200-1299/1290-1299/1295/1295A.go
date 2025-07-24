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
		var n int
		fmt.Fscan(in, &n)
		if n%2 == 1 {
			fmt.Fprint(out, "7")
			n -= 3
		}
		for i := 0; i < n; i += 2 {
			fmt.Fprint(out, "1")
		}
		fmt.Fprintln(out)
	}
}
