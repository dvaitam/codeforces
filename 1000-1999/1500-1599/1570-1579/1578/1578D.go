package main

import (
	"bufio"
	"fmt"
	"os"
)

func dragonCurve(x, y int64) (int, int64) {
	bit := int64(1)
	pos := int64(0)
	for !(x >= -1 && x <= 0 && y >= -1 && y <= 0) {
		bit <<= 1
		if ((x ^ y ^ ((x ^ y) >> 1)) & 1) != 0 {
			pos = bit - 1 - pos
		}
		if ((x ^ y) >> 1 & 1) != 0 {
			nx := (x >> 1) + ((y + 1) >> 1)
			ny := ((y + 1) >> 1) - (x >> 1) - 1
			x, y = nx, ny
		} else {
			nx := ((x + 1) >> 1) + (y >> 1)
			ny := (y >> 1) - ((x + 1) >> 1)
			x, y = nx, ny
		}
	}
	curve := int((x&1 ^ y&3) + 1)
	return curve, pos + 1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	for i := 0; i < n; i++ {
		var x, y int64
		fmt.Fscan(in, &x, &y)
		curve, pos := dragonCurve(x, y)
		fmt.Fprintf(out, "%d %d\n", curve, pos)
	}
}
