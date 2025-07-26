package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

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
			fmt.Fprintln(out, 0)
			continue
		}
		half := n / 2
		factorial := 1
		for i := 2; i <= half; i++ {
			factorial = factorial * i % mod
		}
		ans := factorial * factorial % mod
		fmt.Fprintln(out, ans)
	}
}
