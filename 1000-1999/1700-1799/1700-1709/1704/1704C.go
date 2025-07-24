package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int64
		fmt.Fscan(reader, &n, &m)
		arr := make([]int64, m)
		for i := int64(0); i < m; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })

		gaps := make([]int64, 0, m)
		for i := int64(0); i < m-1; i++ {
			gap := arr[i+1] - arr[i] - 1
			if gap > 0 {
				gaps = append(gaps, gap)
			}
		}
		lastGap := n - arr[m-1] + arr[0] - 1
		if lastGap > 0 {
			gaps = append(gaps, lastGap)
		}
		sort.Slice(gaps, func(i, j int) bool { return gaps[i] > gaps[j] })

		infected := m
		days := int64(0)
		for _, g := range gaps {
			g -= 2 * days
			if g <= 0 {
				continue
			}
			if g == 1 {
				infected += 1
				days += 1
			} else {
				infected += g - 1
				days += 2
			}
		}
		fmt.Fprintln(writer, infected)
	}
}
