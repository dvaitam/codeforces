package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		var a, b string
		fmt.Fscan(in, &a)
		fmt.Fscan(in, &b)

		// collect unique letters in a
		letterIndex := make(map[byte]int)
		letters := make([]byte, 0, 10)
		for i := 0; i < n; i++ {
			c := a[i]
			if _, ok := letterIndex[c]; !ok {
				letterIndex[c] = len(letters)
				letters = append(letters, c)
			}
		}
		m := len(letters)
		if k >= m {
			// we can change all letters
			total := int64(n) * int64(n+1) / 2
			fmt.Fprintln(out, total)
			continue
		}
		ans := int64(0)
		maxMask := 1 << m
		for mask := 0; mask < maxMask; mask++ {
			if bits.OnesCount(uint(mask)) > k {
				continue
			}
			cur := int64(0)
			curLen := int64(0)
			for i := 0; i < n; i++ {
				if a[i] == b[i] {
					curLen++
				} else {
					idx := letterIndex[a[i]]
					if (mask>>idx)&1 == 1 {
						curLen++
					} else {
						cur += curLen * (curLen + 1) / 2
						curLen = 0
					}
				}
			}
			cur += curLen * (curLen + 1) / 2
			if cur > ans {
				ans = cur
			}
		}
		fmt.Fprintln(out, ans)
	}
}
