package main

import (
	"bufio"
	"fmt"
	"os"
)

const limit uint64 = 1 << 61

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var b, c, d uint64
		fmt.Fscan(reader, &b, &c, &d)

		var a uint64
		ok := true

		for bit := 0; bit <= 61; bit++ {
			mask := uint64(1) << bit
			bb := (b >> bit) & 1
			cc := (c >> bit) & 1
			dd := (d >> bit) & 1

			switch (bb << 1) | cc {
			case 0: // b=0, c=0
				if dd == 1 {
					a |= mask
				}
			case 1: // b=0, c=1
				if dd == 1 {
					ok = false
				}
				// choose ai = 0
			case 2: // b=1, c=0
				if dd == 0 {
					ok = false
				}
				// ai free; choose 0
			case 3: // b=1, c=1
				if dd == 1 {
					// need ai = 0 -> nothing
				} else {
					// need ai = 1
					a |= mask
				}
			}
			if !ok {
				break
			}
		}

		if !ok || a > limit || ((a|b)-(a&c)) != d {
			fmt.Fprintln(writer, -1)
		} else {
			fmt.Fprintln(writer, a)
		}
	}
}
