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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	if n <= 5 {
		fmt.Fprintln(out, -1)
	} else {
		tmp := n - 2
		fmt.Fprintln(out, "1 2")
		tot := tmp / 2
		for i := 1; i <= tot; i++ {
			fmt.Fprintf(out, "1 %d\n", i+2)
		}
		for i := 1; i <= tmp-tot; i++ {
			fmt.Fprintf(out, "2 %d\n", i+tot+2)
		}
	}

	for i := 2; i <= n; i++ {
		fmt.Fprintf(out, "1 %d\n", i)
	}
}
