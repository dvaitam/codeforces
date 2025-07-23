package main

import (
	"bufio"
	"fmt"
	"os"
)

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	coords := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &coords[i])
	}

	for i := 0; i < n; i++ {
		var minDist, maxDist int64
		if i == 0 {
			minDist = absInt64(coords[1] - coords[0])
			maxDist = absInt64(coords[n-1] - coords[0])
		} else if i == n-1 {
			minDist = absInt64(coords[n-1] - coords[n-2])
			maxDist = absInt64(coords[n-1] - coords[0])
		} else {
			left := absInt64(coords[i] - coords[i-1])
			right := absInt64(coords[i+1] - coords[i])
			if left < right {
				minDist = left
			} else {
				minDist = right
			}
			distToFirst := absInt64(coords[i] - coords[0])
			distToLast := absInt64(coords[n-1] - coords[i])
			if distToFirst > distToLast {
				maxDist = distToFirst
			} else {
				maxDist = distToLast
			}
		}
		fmt.Fprintf(writer, "%d %d\n", minDist, maxDist)
	}
}
