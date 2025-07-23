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

	var n int
	fmt.Fscan(in, &n)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	minVal := arr[0]
	for i := 1; i < n; i++ {
		if arr[i] < minVal {
			minVal = arr[i]
		}
	}

	prev := -1
	ans := n
	for i := 0; i < n; i++ {
		if arr[i] == minVal {
			if prev != -1 {
				if i-prev < ans {
					ans = i - prev
				}
			}
			prev = i
		}
	}

	fmt.Fprintln(out, ans)
}
