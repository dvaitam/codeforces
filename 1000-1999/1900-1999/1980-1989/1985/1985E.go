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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var x, y, z, k int64
		fmt.Fscan(in, &x, &y, &z, &k)
		var maxWays int64
		for i := int64(1); i <= x; i++ {
			for j := int64(1); j <= y; j++ {
				if k%(i*j) != 0 {
					continue
				}
				p := k / (i * j)
				if p <= z {
					ways := (x - i + 1) * (y - j + 1) * (z - p + 1)
					if ways > maxWays {
						maxWays = ways
					}
				}
			}
		}
		fmt.Fprintln(out, maxWays)
	}
}
