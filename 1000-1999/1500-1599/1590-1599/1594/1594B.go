package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int64
		fmt.Fscan(reader, &n, &k)

		ans := int64(0)
		pow := int64(1)
		for k > 0 {
			if k&1 == 1 {
				ans = (ans + pow) % mod
			}
			pow = (pow * n) % mod
			k >>= 1
		}

		fmt.Fprintln(writer, ans)
	}
}
