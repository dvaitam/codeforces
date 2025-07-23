package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MAX int64 = 1e18

// generate recursively enumerates all products using primes starting at idx
func generate(pr []int64, idx int, cur int64, res *[]int64) {
	if idx == len(pr) {
		*res = append(*res, cur)
		return
	}
	p := pr[idx]
	val := cur
	for val <= MAX {
		generate(pr, idx+1, val, res)
		if val > MAX/p {
			break
		}
		val *= p
	}
}

// countLeq returns how many products a*b are <= x
func countLeq(a, b []int64, x int64) int64 {
	j := len(b) - 1
	var cnt int64
	for _, v := range a {
		for j >= 0 && v > x/b[j] {
			j--
		}
		if j < 0 {
			break
		}
		cnt += int64(j + 1)
	}
	return cnt
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	primes := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &primes[i])
	}
	var k int64
	fmt.Fscan(reader, &k)

	m := n / 2
	leftPr := primes[:m]
	rightPr := primes[m:]

	left := make([]int64, 0)
	right := make([]int64, 0)
	generate(leftPr, 0, 1, &left)
	generate(rightPr, 0, 1, &right)
	sort.Slice(left, func(i, j int) bool { return left[i] < left[j] })
	sort.Slice(right, func(i, j int) bool { return right[i] < right[j] })

	l, r := int64(1), MAX
	for l < r {
		mid := (l + r) / 2
		if countLeq(left, right, mid) >= k {
			r = mid
		} else {
			l = mid + 1
		}
	}
	fmt.Fprintln(writer, l)
}
