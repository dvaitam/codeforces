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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		type sofa struct {
			price int
			idx   int
		}
		arr := make([]sofa, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i].price)
			arr[i].idx = i + 1
		}
		sort.Slice(arr, func(i, j int) bool {
			return arr[i].price < arr[j].price
		})

		sum := 0
		ans := []int{}
		for i := 1; i < n; i++ {
			sum += arr[i-1].price
			if sum < arr[i].price {
				ans = append(ans, arr[i].idx)
			}
		}
		sort.Ints(ans)
		fmt.Fprintln(out, len(ans))
		if len(ans) > 0 {
			for i, v := range ans {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, v)
			}
			fmt.Fprintln(out)
		} else {
			fmt.Fprintln(out)
		}
	}
}
