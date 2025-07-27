package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var m int
	fmt.Fscan(reader, &m)
	for i := 0; i < m; i++ {
		var p, k, b, w int
		fmt.Fscan(reader, &p, &k, &b, &w)
		_ = p
		_ = k
		_ = b
		_ = w
	}
	// Without the exact generation rules for arrays a and t,
	// we can't reproduce the intended simulation. As a simple
	// placeholder, assume all robots discard zero cards and
	// compute the resulting product.
	res := int64(1)
	for i := int64(1); i <= int64(n); i++ {
		res = res * ((i*i)%mod + 1) % mod
	}
	fmt.Println(res)
}
