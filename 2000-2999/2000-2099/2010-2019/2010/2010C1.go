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

	var t string
	fmt.Fscan(in, &t)
	n := len(t)
	for m := n - 1; m > n/2; m-- {
		k := 2*m - n
		if k <= 0 || k >= m {
			continue
		}
		s := t[:m]
		if s[k:] == t[m:] {
			fmt.Fprintln(out, "YES")
			fmt.Fprintln(out, s)
			return
		}
	}
	fmt.Fprintln(out, "NO")
}
