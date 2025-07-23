package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	left := make([][2]int, 0, n)
	right := make([][2]int, 0, n)
	for i := 0; i < n; i++ {
		var x, a int
		fmt.Fscan(reader, &x, &a)
		if x < 0 {
			left = append(left, [2]int{x, a})
		} else {
			right = append(right, [2]int{x, a})
		}
	}
	sort.Slice(left, func(i, j int) bool {
		return left[i][0] > left[j][0]
	})
	sort.Slice(right, func(i, j int) bool {
		return right[i][0] < right[j][0]
	})

	k := len(left)
	if len(right) < k {
		k = len(right)
	}
	sum := 0
	for i := 0; i < k; i++ {
		sum += left[i][1]
		sum += right[i][1]
	}
	if len(left) > len(right) && len(left) > k {
		sum += left[k][1]
	} else if len(right) > len(left) && len(right) > k {
		sum += right[k][1]
	}
	fmt.Println(sum)
}
