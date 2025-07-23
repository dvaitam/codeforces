package main

import (
	"bufio"
	"fmt"
	"os"
)

func countPairs(n, sum int64) int64 {
	if sum < 2 || sum > 2*n {
		return 0
	}
	iMin := int64(1)
	if sum-n > iMin {
		iMin = sum - n
	}
	iMax := n
	if sum-1 < iMax {
		iMax = sum - 1
	}
	maxI := (sum - 1) / 2
	if iMax > maxI {
		iMax = maxI
	}
	if iMin > iMax {
		return 0
	}
	return iMax - iMin + 1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int64
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	pow10 := [10]int64{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000}
	for k := 9; k >= 0; k-- {
		m := pow10[k]
		r := m - 1
		var total int64
		for s := r; s <= 2*n; s += m {
			if s < 2 {
				continue
			}
			total += countPairs(n, s)
		}
		if total > 0 {
			fmt.Fprintln(out, total)
			return
		}
	}
}
