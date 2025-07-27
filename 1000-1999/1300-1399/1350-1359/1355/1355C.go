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

	var A, B, C, D int
	if _, err := fmt.Fscan(in, &A, &B, &C, &D); err != nil {
		return
	}

	maxS := B + C
	diff := make([]int64, maxS+2)

	for x := A; x <= B; x++ {
		start := x + B
		end := x + C
		diff[start]++
		diff[end+1]--
	}

	var ans int64
	var pairs int64
	for s := A + B; s <= maxS; s++ {
		pairs += diff[s]
		t := s - 1
		if t >= C {
			zMax := t
			if zMax > D {
				zMax = D
			}
			if zMax >= C {
				zCount := int64(zMax - C + 1)
				ans += pairs * zCount
			}
		}
	}

	fmt.Fprintln(out, ans)
}
