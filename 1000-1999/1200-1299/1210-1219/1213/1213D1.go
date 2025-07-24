package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	costs := make(map[int][]int)
	for _, v := range a {
		x := v
		ops := 0
		for {
			costs[x] = append(costs[x], ops)
			if x == 0 {
				break
			}
			x /= 2
			ops++
		}
	}

	ans := int(1<<31 - 1)
	for _, arr := range costs {
		if len(arr) < k {
			continue
		}
		sort.Ints(arr)
		sum := 0
		for i := 0; i < k; i++ {
			sum += arr[i]
		}
		if sum < ans {
			ans = sum
		}
	}

	fmt.Fprintln(out, ans)
}
