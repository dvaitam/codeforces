package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		cnt := make(map[int64]int)
		arr := make([]int64, n)
		minVal := int64(1<<63 - 1)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			cnt[arr[i]]++
			if arr[i] < minVal {
				minVal = arr[i]
			}
		}

		if cnt[minVal] >= (n+1)/2 {
			fmt.Fprintln(out, "Yes")
			continue
		}

		divCount := make(map[int64]int)
		for i := 0; i < n; i++ {
			g := gcd(arr[i], minVal)
			divCount[g]++
		}

		divisible := 0
		for _, v := range divCount {
			if v == n {
				divisible++
			}
		}

		if divisible > 0 {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
