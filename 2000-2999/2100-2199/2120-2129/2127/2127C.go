package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k) // k is unused; one optimal pair suffices

		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		base := int64(0)
		segments := make([][2]int64, n)
		for i := 0; i < n; i++ {
			base += abs64(a[i] - b[i])
			if a[i] < b[i] {
				segments[i] = [2]int64{a[i], b[i]}
			} else {
				segments[i] = [2]int64{b[i], a[i]}
			}
		}

		sort.Slice(segments, func(i, j int) bool {
			if segments[i][0] == segments[j][0] {
				return segments[i][1] < segments[j][1]
			}
			return segments[i][0] < segments[j][0]
		})

		minDelta := int64(0)
		maxR := segments[0][1]
		overlap := false
		minGap := int64(1 << 62)
		for i := 1; i < n; i++ {
			l, r := segments[i][0], segments[i][1]
			if l <= maxR {
				overlap = true
				break
			}
			if l-maxR < minGap {
				minGap = l - maxR
			}
			if r > maxR {
				maxR = r
			}
		}

		if !overlap {
			minDelta = 2 * minGap
		}

		fmt.Fprintln(out, base+minDelta)
	}
}
