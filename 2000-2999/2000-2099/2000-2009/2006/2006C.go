package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

const (
	maxDepth = 20
)

var factors []int

func init() {
	maxVal := 1000
	factors = make([]int, maxVal+1)
	for i := 2; i <= maxVal; i++ {
		for j := 2; j*j <= i; j++ {
			if i%j == 0 {
				factors[i] = j
				break
			}
		}
		if factors[i] == 0 {
			factors[i] = 1
		}
	}
}

func isPowerOfTwo(x int) bool {
	return x > 0 && (x&(x-1)) == 0
}

func closureInterval(minVal, maxVal int64, nodes map[int64]struct{}, depth int) bool {
	size := len(nodes)
	if size > maxDepth {
		return maxVal-minVal+1 <= int64(size)*2
	}
	for iter := 0; iter < maxDepth; iter++ {
		added := false
		cur := make([]int64, 0, len(nodes))
		for v := range nodes {
			cur = append(cur, v)
		}
		sort.Slice(cur, func(i, j int) bool { return cur[i] < cur[j] })
		for i := 0; i < len(cur); i++ {
			for j := i + 1; j < len(cur); j++ {
				if cur[i] > math.MaxInt64-cur[j] {
					continue
				}
				sum := cur[i] + cur[j]
				if sum&1 == 0 {
					mid := sum / 2
					if mid > 0 {
						if _, ok := nodes[mid]; !ok {
							nodes[mid] = struct{}{}
							if mid < minVal {
								minVal = mid
							}
							if mid > maxVal {
								maxVal = mid
							}
							added = true
						}
					}
				}
			}
		}
		if !added {
			break
		}
		if int64(len(nodes)) >= maxVal-minVal+1 {
			return true
		}
	}
	return int64(len(nodes)) >= maxVal-minVal+1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		ans := int64(0)
		for l := 0; l < n; l++ {
			nodes := make(map[int64]struct{})
			minVal := int64(1 << 60)
			maxVal := int64(-1 << 60)
			g := int64(0)
			for r := l; r < n && r-l < 40; r++ {
				val := a[r]
				nodes[val] = struct{}{}
				if val < minVal {
					minVal = val
				}
				if val > maxVal {
					maxVal = val
				}
				g = gcd(g, val)
				if maxVal-minVal+1 <= int64(len(nodes))*2 {
					if closureInterval(minVal, maxVal, mapCopy(nodes), 0) {
						ans++
					}
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}

func mapCopy(src map[int64]struct{}) map[int64]struct{} {
	dst := make(map[int64]struct{}, len(src))
	for k := range src {
		dst[k] = struct{}{}
	}
	return dst
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
