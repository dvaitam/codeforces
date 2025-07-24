package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1_000_000_007

func modPow(a, n int64) int64 {
	res := int64(1)
	for n > 0 {
		if n&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		n >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int64
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	val := modPow(2, m) - 1
	if val < 0 {
		val += mod
	}
	ans := modPow(val, n)
	fmt.Fprintln(writer, ans)
}
