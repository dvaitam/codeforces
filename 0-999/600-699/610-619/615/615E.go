package main

import (
	"bufio"
	"fmt"
	"os"
)

func spiralPos(n int64) (int64, int64) {
	if n == 0 {
		return 0, 0
	}
	low, high := int64(0), int64(1000000000)
	for low < high {
		mid := (low + high) / 2
		if 3*mid*(mid+1) >= n {
			high = mid
		} else {
			low = mid + 1
		}
	}
	r := low
	prev := 3 * (r - 1) * r
	k := n - prev
	x, y := r, int64(0)
	dirs := [6][2]int64{{-1, 1}, {-1, 0}, {0, -1}, {1, -1}, {1, 0}, {0, 1}}
	for i := 0; i < 6 && k > 0; i++ {
		step := r
		if k < step {
			step = k
		}
		x += dirs[i][0] * step
		y += dirs[i][1] * step
		k -= step
	}
	return x, y
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int64
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	x, y := spiralPos(n)
	fmt.Fprintf(writer, "%d %d\n", x, y)
}
