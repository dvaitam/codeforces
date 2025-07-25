package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MaxV = 100000

func totients(n int) []int64 {
	phi := make([]int64, n+1)
	for i := 0; i <= n; i++ {
		phi[i] = int64(i)
	}
	for i := 2; i <= n; i++ {
		if phi[i] == int64(i) {
			for j := i; j <= n; j += i {
				phi[j] -= phi[j] / int64(i)
			}
		}
	}
	return phi
}

func divisors(x int) []int {
	res := []int{}
	for d := 1; d*d <= x; d++ {
		if x%d == 0 {
			res = append(res, d)
			if d*d != x {
				res = append(res, x/d)
			}
		}
	}
	return res
}

func main() {
	phi := totients(MaxV)
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		sort.Ints(arr)
		freq := make([]int64, MaxV+1)
		ans := int64(0)
		for j, val := range arr {
			ds := divisors(val)
			total := int64(0)
			for _, d := range ds {
				total += phi[d] * freq[d]
			}
			ans += total * int64(n-j-1)
			for _, d := range ds {
				freq[d]++
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
