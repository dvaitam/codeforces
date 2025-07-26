package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func powmod(a int64, b int64) int64 {
	a %= mod
	res := int64(1)
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
		var n int
		fmt.Fscan(reader, &n)
		edges := make(map[[2]int]struct{}, n-1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			if u > v {
				u, v = v, u
			}
			edges[[2]int{u, v}] = struct{}{}
		}
		common := 0
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			if u > v {
				u, v = v, u
			}
			if _, ok := edges[[2]int{u, v}]; ok {
				common++
			}
		}
		k := int64(n - 1 - common)
		ans := powmod(2, k)
		fmt.Fprintln(writer, ans)
	}
}
