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
		var n int
		var k int64
		fmt.Fscan(reader, &n, &k)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		prefix := make([]int64, n+1)
		for i := 0; i < n; i++ {
			prefix[i+1] = prefix[i] + a[i]
		}
		total := prefix[n]
		if total <= k {
			fmt.Fprintln(writer, 0)
			continue
		}
		ans := int64(1 << 60)
		a0 := a[0]
		for m := 0; m < n; m++ {
			sumLast := prefix[n] - prefix[n-m]
			// sum after setting m largest elements equal to a0
			s1 := prefix[n] - sumLast + int64(m)*a0
			diff := s1 - k
			var x int64
			if diff <= 0 {
				x = 0
			} else {
				x = (diff + int64(m)) / int64(m+1)
			}
			if x < 0 {
				x = 0
			}
			steps := int64(m) + x
			if steps < ans {
				ans = steps
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
