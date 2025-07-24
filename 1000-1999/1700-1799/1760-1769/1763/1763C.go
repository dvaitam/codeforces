package main

import (
	"bufio"
	"fmt"
	"os"
)

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		if n == 2 {
			ans := 2 * absInt64(a[0]-a[1])
			fmt.Fprintln(writer, ans)
			continue
		}
		if n == 3 {
			maxVal, minVal := a[0], a[0]
			for i := 1; i < 3; i++ {
				if a[i] > maxVal {
					maxVal = a[i]
				}
				if a[i] < minVal {
					minVal = a[i]
				}
			}
			boundaryMax := a[0]
			if a[2] > boundaryMax {
				boundaryMax = a[2]
			}
			diff := maxVal - minVal
			v := boundaryMax
			if diff > v {
				v = diff
			}
			ans := int64(3) * v
			fmt.Fprintln(writer, ans)
			continue
		}
		var maxVal int64 = a[0]
		for i := 1; i < n; i++ {
			if a[i] > maxVal {
				maxVal = a[i]
			}
		}
		ans := int64(n) * maxVal
		fmt.Fprintln(writer, ans)
	}
}
