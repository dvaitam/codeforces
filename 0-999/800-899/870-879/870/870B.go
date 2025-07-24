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
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	minVal, maxVal := a[0], a[0]
	for i := 1; i < n; i++ {
		if a[i] < minVal {
			minVal = a[i]
		}
		if a[i] > maxVal {
			maxVal = a[i]
		}
	}

	var ans int
	switch {
	case k == 1:
		ans = minVal
	case k >= 3:
		ans = maxVal
	default: // k == 2
		if a[0] > a[n-1] {
			ans = a[0]
		} else {
			ans = a[n-1]
		}
	}

	fmt.Fprintln(writer, ans)
}
