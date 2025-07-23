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

	var n, a, b, c int64
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	fmt.Fscan(in, &a)
	fmt.Fscan(in, &b)
	fmt.Fscan(in, &c)

	var ans int64
	if n >= b && b-c <= a {
		k := (n-b)/(b-c) + 1
		ans += k
		n -= k * (b - c)
	}
	ans += n / a
	fmt.Fprintln(out, ans)
}
