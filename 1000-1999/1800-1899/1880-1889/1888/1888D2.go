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

	product := int64(1)
	for i := 0; i < n; i++ {
		var v int64
		if _, err := fmt.Fscan(in, &v); err != nil {
			return
		}
		v %= mod
		if v < 0 {
			v += mod
		}
		product = (product * v) % mod
	}

	fmt.Fprintln(out, product%mod)
}
