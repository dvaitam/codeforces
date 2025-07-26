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

	var input int
	if _, err := fmt.Fscan(in, &input); err != nil {
		return
	}
	n := input / 1000
	mod := input % 1000

	res := 1 % mod
	for i := n; i > 0; i -= 2 {
		res = (res * (i % mod)) % mod
	}
	fmt.Fprintln(out, res)
}
