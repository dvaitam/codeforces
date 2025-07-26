package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		k := n * m
		nums := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &nums[i])
		}
		sort.Ints(nums)
		mn := nums[0]
		mn2 := nums[1]
		mx := nums[k-1]
		mx2 := nums[k-2]

		s1 := int64(n-1)*int64(mx-mn2) + int64(n)*int64(m-1)*int64(mx-mn)
		s2 := int64(m-1)*int64(mx-mn2) + int64(m)*int64(n-1)*int64(mx-mn)
		s3 := int64(n-1)*int64(mx2-mn) + int64(n)*int64(m-1)*int64(mx-mn)
		s4 := int64(m-1)*int64(mx2-mn) + int64(m)*int64(n-1)*int64(mx-mn)

		ans := s1
		if s2 > ans {
			ans = s2
		}
		if s3 > ans {
			ans = s3
		}
		if s4 > ans {
			ans = s4
		}

		fmt.Fprintln(writer, ans)
	}
}
