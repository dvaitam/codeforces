package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	const mask uint32 = (1 << 26) - 1

	cnt := make(map[uint64]int)
	var ans int64

	for ; n > 0; n-- {
		var s string
		fmt.Fscan(in, &s)
		var pres, parity uint32
		for i := 0; i < len(s); i++ {
			bit := uint32(1) << (s[i] - 'a')
			pres |= bit
			parity ^= bit
		}
		abs := mask ^ pres
		lp := uint32(len(s) & 1)
		for l := uint32(0); l < 26; l++ {
			if abs&(1<<l) == 0 {
				continue
			}
			target := ^(parity ^ (1 << l)) & mask
			key := (uint64(target) << 6) | (uint64(l) << 1) | uint64(lp^1)
			ans += int64(cnt[key])
		}
		for l := uint32(0); l < 26; l++ {
			if abs&(1<<l) == 0 {
				continue
			}
			key := (uint64(parity) << 6) | (uint64(l) << 1) | uint64(lp)
			cnt[key]++
		}
	}
	fmt.Println(ans)
}
