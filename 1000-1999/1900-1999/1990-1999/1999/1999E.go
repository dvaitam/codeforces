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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var l, r int64
		fmt.Fscan(in, &l, &r)

		count := r - l + 1
		const mod = int64(998244353)
		k := count % mod
		n := count % mod
		tmp := (k * ((n - 1 + mod) % mod)) % mod
		tmp = (tmp * 499122177) % mod // inverse of 2 mod
		result := (count*k%mod - tmp + mod) % mod

		fmt.Fprintln(out, result)
	}
}
