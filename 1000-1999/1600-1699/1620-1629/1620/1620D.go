package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		const INF int64 = 1 << 60
		best := INF
		for c1 := 0; c1 <= 2; c1++ {
			for c2 := 0; c2 <= 2; c2++ {
				max3 := int64(0)
				feasible := true
				for _, x := range a {
					min3 := INF
					for i1 := 0; i1 <= c1; i1++ {
						for i2 := 0; i2 <= c2; i2++ {
							rem := x - int64(i1) - 2*int64(i2)
							if rem >= 0 && rem%3 == 0 {
								tmp := rem / 3
								if tmp < min3 {
									min3 = tmp
								}
							}
						}
					}
					if min3 == INF {
						feasible = false
						break
					}
					if min3 > max3 {
						max3 = min3
					}
				}
				if feasible {
					total := int64(c1+c2) + max3
					if total < best {
						best = total
					}
				}
			}
		}
		fmt.Fprintln(writer, best)
	}
}
