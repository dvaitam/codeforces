package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	pow10 := make([]int64, n+1)
	pow10[0] = 1
	for i := 1; i <= n; i++ {
		pow10[i] = pow10[i-1] * 10 % mod
	}

	for i := 1; i <= n; i++ {
		var res int64
		if i == n {
			res = 10
		} else {
			part1 := int64(2*10*9) * pow10[n-i-1] % mod
			res = part1
			if n-i-1 > 0 {
				part2 := int64(n-i-1) * 10 % mod
				part2 = part2 * 9 % mod
				part2 = part2 * 9 % mod
				part2 = part2 * pow10[n-i-2] % mod
				res = (res + part2) % mod
			}
		}
		if i > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, res)
	}
	fmt.Fprintln(out)
}
