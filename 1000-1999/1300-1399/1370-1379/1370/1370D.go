package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	fmt.Fscan(reader, &n, &k)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	low, high := 1, int(1e9)
	answer := high

	check := func(x int) bool {
		// pattern 1: positions 1,3,... must be <= x
		cnt := 0
		for _, v := range arr {
			if cnt%2 == 0 { // expecting constrained position
				if v <= x {
					cnt++
				}
			} else {
				cnt++
			}
			if cnt >= k {
				return true
			}
		}
		// pattern 2: positions 2,4,... must be <= x
		cnt = 0
		for _, v := range arr {
			if cnt%2 == 1 { // even position requires <= x
				if v <= x {
					cnt++
				}
			} else {
				cnt++
			}
			if cnt >= k {
				return true
			}
		}
		return false
	}

	for low <= high {
		mid := (low + high) / 2
		if check(mid) {
			answer = mid
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	fmt.Fprintln(writer, answer)
}
