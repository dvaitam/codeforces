package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

func canReach(prices []int64, x, a, y, b int64, k int64, m int) bool {
	l := lcm(a, b)
	mm := int64(m)
	cBoth := mm / l
	cA := mm/a - cBoth
	cB := mm/b - cBoth
	idx := 0
	var sum int64
	for i := int64(0); i < cBoth && idx < len(prices); i++ {
		sum += prices[idx] * (x + y)
		idx++
	}
	for i := int64(0); i < cA && idx < len(prices); i++ {
		sum += prices[idx] * x
		idx++
	}
	for i := int64(0); i < cB && idx < len(prices); i++ {
		sum += prices[idx] * y
		idx++
	}
	return sum >= k
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var n int
		fmt.Fscan(reader, &n)
		prices := make([]int64, n)
		for i := 0; i < n; i++ {
			var p int64
			fmt.Fscan(reader, &p)
			prices[i] = p / 100
		}
		sort.Slice(prices, func(i, j int) bool { return prices[i] > prices[j] })

		var x, a int64
		fmt.Fscan(reader, &x, &a)
		var y, b int64
		fmt.Fscan(reader, &y, &b)
		if x < y {
			x, y = y, x
			a, b = b, a
		}
		var k int64
		fmt.Fscan(reader, &k)
		left, right := 1, n
		ans := -1
		for left <= right {
			mid := (left + right) / 2
			if canReach(prices, x, a, y, b, k, mid) {
				ans = mid
				right = mid - 1
			} else {
				left = mid + 1
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
