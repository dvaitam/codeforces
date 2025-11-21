package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type point struct {
	x int64
	y int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var k int
		fmt.Fscan(in, &k)

		if k == 0 {
			fmt.Fprintln(out, 1)
			fmt.Fprintln(out, "0 0")
			continue
		}

		remainingPairs := k
		totalPoints := 0
		groupSizes := make([]int, 0)

		for remainingPairs > 0 {
			remainingSlots := 500 - totalPoints
			if remainingSlots < 2 {
				remainingSlots = 2
			}

			estimate := int((1 + math.Sqrt(float64(1+8*remainingPairs))) / 2)
			if estimate > remainingSlots {
				estimate = remainingSlots
			}
			if estimate < 2 {
				estimate = 2
			}
			for estimate > remainingSlots || combination(estimate) > remainingPairs {
				estimate--
			}

			groupSizes = append(groupSizes, estimate)
			remainingPairs -= combination(estimate)
			totalPoints += estimate
		}

		points := make([]point, 0, totalPoints)
		curY := int64(1)
		for idx, size := range groupSizes {
			x := int64(idx * 2)
			for j := 0; j < size; j++ {
				points = append(points, point{x: x, y: curY})
				curY++
			}
		}

		fmt.Fprintln(out, len(points))
		for _, p := range points {
			fmt.Fprintf(out, "%d %d\n", p.x, p.y)
		}
	}
}

func combination(x int) int {
	return x * (x - 1) / 2
}
