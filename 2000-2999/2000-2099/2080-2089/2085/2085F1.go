package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
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
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		freq := make([]int, k+1)
		distinct := 0
		ans := math.MaxInt32
		left := 0

		for right := 0; right < n; right++ {
			val := a[right]
			freq[val]++
			if freq[val] == 1 {
				distinct++
			}

			for distinct == k {
				if cur := right - left + 1; cur < ans {
					ans = cur
				}
				leftVal := a[left]
				freq[leftVal]--
				if freq[leftVal] == 0 {
					distinct--
				}
				left++
			}
		}

		fmt.Fprintln(out, ans-k)
	}
}
