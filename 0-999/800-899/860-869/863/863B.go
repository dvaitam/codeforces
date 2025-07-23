package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	w := make([]int, 2*n)
	for i := 0; i < 2*n; i++ {
		fmt.Fscan(in, &w[i])
	}
	sort.Ints(w)
	ans := int(^uint(0) >> 1)
	for i := 0; i < 2*n; i++ {
		for j := i + 1; j < 2*n; j++ {
			arr := make([]int, 0, 2*n-2)
			for k := 0; k < 2*n; k++ {
				if k == i || k == j {
					continue
				}
				arr = append(arr, w[k])
			}
			sum := 0
			for k := 0; k < len(arr); k += 2 {
				sum += arr[k+1] - arr[k]
			}
			if sum < ans {
				ans = sum
			}
		}
	}
	fmt.Println(ans)
}
