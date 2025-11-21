package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// solveSingle returns the minimal number of presses needed to guarantee k cans.
func solveSingle(a []int64, k int64) int64 {
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })

	var presses int64
	prev := int64(0) // highest fully processed level
	n := len(a)      // total slots
	processed := 0   // number of slots with capacity <= prev

	for processed < n && k > 0 {
		nextVal := a[processed]
		height := nextVal - prev
		if height > 0 {
			available := int64(n - processed)
			capacity := available * height // successes available in this block
			if capacity >= k {
				presses += k
				k = 0
				break
			}
			presses += capacity
			k -= capacity
		}

		// Need to move to the next level: pay penalty for slots that end at nextVal.
		cntSame := 0
		for processed < n && a[processed] == nextVal {
			processed++
			cntSame++
		}
		presses += int64(cntSame) // failed presses when trying to go beyond nextVal
		prev = nextVal
	}

	return presses
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)

		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}

		result := solveSingle(arr, k)
		fmt.Fprintln(out, result)
	}
}
