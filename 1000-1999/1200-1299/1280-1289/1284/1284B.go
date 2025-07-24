package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)

	minVals := make([]int, 0, n)
	maxVals := make([]int, 0, n)
	for i := 0; i < n; i++ {
		var l int
		fmt.Fscan(reader, &l)
		arr := make([]int, l)
		for j := 0; j < l; j++ {
			fmt.Fscan(reader, &arr[j])
		}
		hasAscent := false
		minVal := arr[0]
		maxVal := arr[0]
		for j := 1; j < l; j++ {
			if arr[j] > arr[j-1] {
				hasAscent = true
			}
			if arr[j] < minVal {
				minVal = arr[j]
			}
			if arr[j] > maxVal {
				maxVal = arr[j]
			}
		}
		if !hasAscent {
			minVals = append(minVals, minVal)
			maxVals = append(maxVals, maxVal)
		}
	}

	total := int64(n) * int64(n)
	sort.Ints(maxVals)
	var bad int64
	for _, mn := range minVals {
		cnt := sort.SearchInts(maxVals, mn+1)
		bad += int64(cnt)
	}
	fmt.Println(total - bad)
}
