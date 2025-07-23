package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, h, w int
	if _, err := fmt.Fscan(in, &n, &h, &w); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)
	remW := int64(w)
	remH := int64(h)
	dx, dy := 0, 0
	minX, maxX := 0, 0
	minY, maxY := 0, 0
	ans := int64(0)
	step := 0
	idx := 0
	for remW > 0 && remH > 0 {
		ch := s[idx]
		idx++
		if idx == n {
			idx = 0
		}
		step++
		switch ch {
		case 'L':
			dx--
		case 'R':
			dx++
		case 'U':
			dy--
		case 'D':
			dy++
		}
		changed := false
		if dx < minX {
			minX = dx
			remW--
			ans = (ans + remH*int64(step)) % MOD
			changed = true
		}
		if dx > maxX && remW > 0 {
			maxX = dx
			remW--
			ans = (ans + remH*int64(step)) % MOD
			changed = true
		}
		if dy < minY && remH > 0 {
			minY = dy
			remH--
			ans = (ans + remW*int64(step)) % MOD
			changed = true
		}
		if dy > maxY && remH > 0 {
			maxY = dy
			remH--
			ans = (ans + remW*int64(step)) % MOD
			changed = true
		}
		if !changed && idx == 0 && dx == 0 && dy == 0 {
			fmt.Println(-1)
			return
		}
	}
	fmt.Println(ans % MOD)
}
