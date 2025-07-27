package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	// read assignment permutation but it does not affect the result
	for i := 0; i < n; i++ {
		var tmp int
		fmt.Fscan(in, &tmp)
	}

	if n == 1 {
		fmt.Fprintln(out, 1%mod)
	} else {
		fmt.Fprintln(out, 0)
	}
}
