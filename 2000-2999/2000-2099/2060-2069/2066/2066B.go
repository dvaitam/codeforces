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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		zeroCnt := 0
		firstZero := -1
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			if arr[i] == 0 {
				if firstZero == -1 {
					firstZero = i
				}
				zeroCnt++
			}
		}

		if zeroCnt == 0 {
			// 0 is missing, so mex of any suffix is 0 and the whole sequence works.
			fmt.Fprintln(out, n)
			continue
		}

		// Keep only the earliest zero; others are discarded.
		ans := n - (zeroCnt - 1)

		present := make([]bool, n+2)
		mex := 1 // suffix already contains a zero after we pass its position
		hasZeroInSuffix := false
		bad := false

		for i := n - 1; i >= 0; i-- {
			if hasZeroInSuffix && mex > arr[i] {
				bad = true
				break
			}

			if arr[i] <= n {
				if !present[arr[i]] {
					present[arr[i]] = true
					for present[mex] {
						mex++
					}
				}
			}

			if i == firstZero {
				hasZeroInSuffix = true
			}
		}

		if bad {
			ans-- // remove the kept zero as well
		}

		fmt.Fprintln(out, ans)
	}
}
