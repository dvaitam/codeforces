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

	var n, q, m int
	if _, err := fmt.Fscan(in, &n, &q, &m); err != nil {
		return
	}
	for i := 0; i < n; i++ {
		var tmp int
		fmt.Fscan(in, &tmp)
	}
	for i := 0; i < n; i++ {
		var tmp int
		fmt.Fscan(in, &tmp)
	}
	for i := 0; i < n; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
	}
	for i := 0; i < m; i++ {
		var u, p int
		if _, err := fmt.Fscan(in, &u, &p); err != nil {
			return
		}
		fmt.Fprintln(out, 0)
	}
}
