package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func powMod(a, b int64) int64 {
	res := int64(1)
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, k int64
		fmt.Fscan(reader, &n, &k)
		ans := powMod(n, k)
		fmt.Fprintln(writer, ans)
	}
}
