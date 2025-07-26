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

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	half := k / 2
	total := 0
	for l := 0; l <= n-k; l++ {
		mism := 0
		for d := 0; d < half; d++ {
			if arr[l+d] != arr[l+k-1-d] {
				mism++
			}
		}
		total += mism
	}
	fmt.Fprintln(out, total)
}
