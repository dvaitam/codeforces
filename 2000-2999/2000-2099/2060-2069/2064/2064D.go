package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func msb(x int) int {
	if x == 0 {
		return -1
	}
	return bits.Len(uint(x)) - 1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		w := make([]int, n+1)  // 1-based
		mb := make([]int, n+1) // msb of each weight
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &w[i])
			mb[i] = msb(w[i])
		}
		prefix := make([]int, n+1)
		for i := 1; i <= n; i++ {
			prefix[i] = prefix[i-1] ^ w[i]
		}

		// prev[i][b]: nearest position <= i with msb >= b, 0 if none
		prev := make([][30]int, n+1)
		for b := 0; b < 30; b++ {
			prev[0][b] = 0
		}
		for i := 1; i <= n; i++ {
			for b := 0; b < 30; b++ {
				if mb[i] >= b {
					prev[i][b] = i
				} else {
					prev[i][b] = prev[i-1][b]
				}
			}
		}

		for ; q > 0; q-- {
			var x int
			fmt.Fscan(in, &x)
			pos := n
			cur := x
			eaten := 0
			for pos > 0 && cur > 0 {
				b := msb(cur)
				j := prev[pos][b]
				if j < pos {
					// eat slimes with smaller msb and apply their xor
					segXor := 0
					if j+1 <= pos {
						segXor = prefix[pos] ^ prefix[j]
					}
					cur ^= segXor
					eaten += pos - j
					pos = j
				}
				if pos == 0 {
					break
				}

				if mb[pos] > b {
					break // cannot eat next slime
				}
				if cur < w[pos] {
					break
				}
				cur ^= w[pos]
				eaten++
				pos--
			}
			fmt.Fprintln(out, eaten)
		}
	}
}
