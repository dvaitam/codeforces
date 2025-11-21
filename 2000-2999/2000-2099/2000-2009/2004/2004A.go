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
		var n int
		fmt.Fscan(in, &n)
		points := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &points[i])
		}

		distances := make([]int, n)
		for i := 0; i < n; i++ {
			if i == 0 {
				distances[i] = points[1] - points[0]
			} else if i == n-1 {
				distances[i] = points[n-1] - points[n-2]
			} else {
				prev := points[i] - points[i-1]
				next := points[i+1] - points[i]
				if prev < next {
					distances[i] = prev
				} else {
					distances[i] = next
				}
			}
		}

		left := -1 << 60
		right := 1 << 60
		for i := 0; i < n; i++ {
			if cur := points[i] - distances[i]; cur > left {
				left = cur
			}
			if cur := points[i] + distances[i]; cur < right {
				right = cur
			}
		}

		if left > right {
			fmt.Fprintln(out, "NO")
			continue
		}

		exist := make(map[int]struct{}, n)
		for _, v := range points {
			exist[v] = struct{}{}
		}

		found := false
		for y := left; y <= right; y++ {
			if _, ok := exist[y]; ok {
				continue
			}
			found = true
			break
		}

		if found {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
