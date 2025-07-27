package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		sort.Ints(arr)

		nonPos := make([]int, 0)
		posMin := int(1<<31 - 1)
		for _, v := range arr {
			if v <= 0 {
				nonPos = append(nonPos, v)
			} else if v < posMin {
				posMin = v
			}
		}
		minDiff := int(1<<31 - 1)
		for i := 1; i < len(nonPos); i++ {
			diff := nonPos[i] - nonPos[i-1]
			if diff < minDiff {
				minDiff = diff
			}
		}
		ans := len(nonPos)
		if len(nonPos) <= 1 {
			minDiff = int(1<<31 - 1)
		}
		if posMin != int(1<<31-1) && posMin <= minDiff {
			ans++
		}
		fmt.Fprintln(out, ans)
	}
}
