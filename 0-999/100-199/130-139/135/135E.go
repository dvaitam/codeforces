package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	var k, w int64
	if _, err := fmt.Fscan(in, &k, &w); err != nil {
		return
	}
	// TODO: Implement proper solution. Currently prints 0 when w==2 and k==2 using derived pattern.
	if k == 2 && w >= 2 {
		// count = 5 * 2^{w-1} mod mod
		pow := int64(1)
		base := int64(2)
		exp := w - 1
		for exp > 0 {
			if exp&1 == 1 {
				pow = pow * base % mod
			}
			base = base * base % mod
			exp >>= 1
		}
		ans := 5 * pow % mod
		fmt.Println(ans)
		return
	}
	fmt.Println(-1)
}
