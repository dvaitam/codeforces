package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	counts := [4]int64{}
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		g := ((x/2)%2)*2 + ((y / 2) % 2)
		counts[g]++
	}

	pair := func(a, b int) int {
		if a^b == 3 {
			return 1
		}
		return 0
	}

	var ans int64
	for a := 0; a < 4; a++ {
		for b := a; b < 4; b++ {
			for c := b; c < 4; c++ {
				if (pair(a, b) ^ pair(b, c) ^ pair(c, a)) == 0 {
					ca, cb, cc := counts[a], counts[b], counts[c]
					if a == b && b == c {
						if ca >= 3 {
							ans += ca * (ca - 1) * (ca - 2) / 6
						}
					} else if a == b {
						if ca >= 2 {
							ans += ca * (ca - 1) / 2 * cc
						}
					} else if b == c {
						if cb >= 2 {
							ans += cb * (cb - 1) / 2 * ca
						}
					} else if a == c {
						if ca >= 2 {
							ans += ca * (ca - 1) / 2 * cb
						}
					} else {
						ans += ca * cb * cc
					}
				}
			}
		}
	}

	fmt.Println(ans)
}
