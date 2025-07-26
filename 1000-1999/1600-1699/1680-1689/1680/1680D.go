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

	var n int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	a := make([]int64, n)
	zeroCount := 0
	sumKnown := int64(0)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] == 0 {
			zeroCount++
		} else {
			sumKnown += a[i]
		}
	}

	if abs64(sumKnown) > int64(zeroCount)*k {
		fmt.Fprintln(writer, -1)
		return
	}

	// backward ranges
	low := make([]int64, n+1)
	high := make([]int64, n+1)
	low[n] = 0
	high[n] = 0
	for i := n - 1; i >= 0; i-- {
		if a[i] == 0 {
			low[i] = low[i+1] - k
			high[i] = high[i+1] + k
		} else {
			low[i] = low[i+1] - a[i]
			high[i] = high[i+1] - a[i]
		}
	}

	curL, curR := int64(0), int64(0)
	minPref, maxPref := int64(0), int64(0)
	for i := 1; i <= n; i++ {
		if a[i-1] == 0 {
			curL -= k
			curR += k
		} else {
			curL += a[i-1]
			curR += a[i-1]
		}
		if curL < low[i] {
			curL = low[i]
		}
		if curR > high[i] {
			curR = high[i]
		}
		if curL > curR {
			// should not happen if input is valid
			fmt.Fprintln(writer, -1)
			return
		}
		if curL < minPref {
			minPref = curL
		}
		if curR > maxPref {
			maxPref = curR
		}
	}

	ans := maxPref - minPref + 1
	fmt.Fprintln(writer, ans)
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
