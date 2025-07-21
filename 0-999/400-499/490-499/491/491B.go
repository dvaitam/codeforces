package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF = 1 << 62

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var N, M int64
	fmt.Fscan(reader, &N, &M)
	var C int
	fmt.Fscan(reader, &C)
	var x, y int64
	// Precompute extrema of transformed hotel coordinates
	var maxS1, maxS2 int64 = -INF, -INF
	var minS1, minS2 int64 = INF, INF
	for i := 0; i < C; i++ {
		fmt.Fscan(reader, &x, &y)
		s1 := x + y
		s2 := x - y
		if s1 > maxS1 {
			maxS1 = s1
		}
		if s1 < minS1 {
			minS1 = s1
		}
		if s2 > maxS2 {
			maxS2 = s2
		}
		if s2 < minS2 {
			minS2 = s2
		}
	}
	var H int
	fmt.Fscan(reader, &H)
	bestIdx := 1
	bestD := INF
	for i := 1; i <= H; i++ {
		fmt.Fscan(reader, &x, &y)
		s1 := x + y
		s2 := x - y
		// compute max Manhattan distance to all hotels
		d1 := maxS1 - s1
		if s1-minS1 > d1 {
			d1 = s1 - minS1
		}
		d2 := maxS2 - s2
		if s2-minS2 > d2 {
			d2 = s2 - minS2
		}
		d := d1
		if d2 > d {
			d = d2
		}
		if d < bestD {
			bestD = d
			bestIdx = i
		}
	}
	fmt.Fprintln(writer, bestD)
	fmt.Fprintln(writer, bestIdx)
}
