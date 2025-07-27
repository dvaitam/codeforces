package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for Codeforces problem "Maximum Subsequence Value" described in
// problemE.txt. An optimal subsequence contains at most three elements:
// for any subsequence of length k >= 4, removing elements can only decrease
// the required count of common bits while keeping all previously counted bits.
// Therefore we simply check every subset of up to three distinct numbers and
// take the maximum bitwise OR among them.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	var ans int64
	for i := 0; i < n; i++ {
		if arr[i] > ans {
			ans = arr[i]
		}
		for j := i + 1; j < n; j++ {
			v := arr[i] | arr[j]
			if v > ans {
				ans = v
			}
			for k := j + 1; k < n; k++ {
				v3 := v | arr[k]
				if v3 > ans {
					ans = v3
				}
			}
		}
	}

	fmt.Fprintln(out, ans)
}
