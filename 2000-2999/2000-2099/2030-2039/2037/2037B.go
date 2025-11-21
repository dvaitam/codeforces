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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var k int64
		fmt.Fscan(in, &k)
		arr := make([]int64, k)
		freq := make(map[int64]int)
		for i := int64(0); i < k; i++ {
			fmt.Fscan(in, &arr[i])
			freq[arr[i]]++
		}

		product := k - 2
		nAns, mAns := int64(1), product
		for d := int64(1); d*d <= product; d++ {
			if product%d != 0 {
				continue
			}
			n, m := d, product/d
			if validPair(n, m, freq) {
				nAns, mAns = n, m
				break
			}
			if validPair(m, n, freq) {
				nAns, mAns = m, n
				break
			}
		}
		fmt.Fprintln(out, nAns, mAns)
	}
}

func validPair(n, m int64, freq map[int64]int) bool {
	if n == m {
		return freq[n] >= 2
	}
	return freq[n] >= 1 && freq[m] >= 1
}
