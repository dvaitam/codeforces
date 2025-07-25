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

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	const maxBit = 60
	for ; q > 0; q-- {
		var k int64
		fmt.Fscan(in, &k)
		cur := make([]int64, n)
		copy(cur, a)
		ans := int64(0)
		for bit := maxBit; bit >= 0; bit-- {
			var need int64
			mask := int64(1) << bit
			for i := 0; i < n; i++ {
				if cur[i]&mask == 0 {
					add := mask - (cur[i] & ((mask << 1) - 1))
					need += add
					if need > k {
						break
					}
				}
			}
			if need <= k {
				ans |= mask
				k -= need
				for i := 0; i < n; i++ {
					if cur[i]&mask == 0 {
						add := mask - (cur[i] & ((mask << 1) - 1))
						cur[i] += add
					}
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
