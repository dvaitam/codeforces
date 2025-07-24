package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const maxX = 100000
const INF = int64(1 << 60)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}

	pos := make([]int64, maxX+2) // minimal positive
	neg := make([]int64, maxX+2) // maximal negative
	val := make([]int64, maxX+2)
	for i := 0; i <= maxX+1; i++ {
		pos[i] = INF
		neg[i] = -INF
	}

	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var l, r int
			var k int64
			fmt.Fscan(reader, &l, &r, &k)
			for x := l; x < r && x <= maxX; x++ {
				if k > 0 {
					if pos[x] > k {
						pos[x] = k
					}
				} else {
					if neg[x] < k {
						neg[x] = k
					}
				}
				if pos[x] != INF && neg[x] != -INF {
					val[x] = pos[x] + int64(math.Abs(float64(neg[x])))
				}
			}
		} else {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			var ans int64
			for x := l; x < r && x <= maxX; x++ {
				if pos[x] != INF && neg[x] != -INF {
					ans += val[x]
				}
			}
			fmt.Fprintln(writer, ans)
		}
	}
}
