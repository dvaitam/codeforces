package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func isSquare(x int64) bool {
	r := int64(math.Sqrt(float64(x)))
	for (r+1)*(r+1) <= x {
		r++
	}
	for r*r > x {
		r--
	}
	return r*r == x
}

func isqrt(x int64) int64 {
	r := int64(math.Sqrt(float64(x)))
	for (r+1)*(r+1) <= x {
		r++
	}
	for r*r > x {
		r--
	}
	return r
}

func precompute() []int64 {
	const limit int64 = 1e18
	set := make(map[int64]struct{})
	var base int64
	for base = 2; base*base*base <= limit; base++ {
		val := base * base * base
		for {
			if val > limit {
				break
			}
			if !isSquare(val) {
				set[val] = struct{}{}
			}
			if val > limit/base {
				break
			}
			val *= base
		}
	}
	arr := make([]int64, 0, len(set))
	for v := range set {
		arr = append(arr, v)
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	return arr
}

func countUpTo(x int64, arr []int64) int64 {
	if x <= 0 {
		return 0
	}
	sq := isqrt(x)
	// count numbers in arr <= x
	idx := sort.Search(len(arr), func(i int) bool { return arr[i] > x })
	return sq + int64(idx)
}

func main() {
	others := precompute()
	reader := bufio.NewReader(os.Stdin)
	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i := 0; i < q; i++ {
		var l, r int64
		fmt.Fscan(reader, &l, &r)
		ans := countUpTo(r, others) - countUpTo(l-1, others)
		fmt.Fprintln(writer, ans)
	}
}
