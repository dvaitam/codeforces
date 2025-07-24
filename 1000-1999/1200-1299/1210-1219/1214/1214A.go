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

	var n, d, e int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	fmt.Fscan(in, &d)
	fmt.Fscan(in, &e)
	e *= 5
	ans := n
	for i := 0; i <= n; i += e {
		rem := (n - i) % d
		if rem < ans {
			ans = rem
		}
	}
	fmt.Fprintln(out, ans)
}
