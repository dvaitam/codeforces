package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	maxDepth = 31 // 2^31 is comfortably above the allowed values
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// computeMinCost returns the minimal number of pizzas needed to make some position
// reach value target while keeping the line comfortable.
// Complexity: O(n * maxDepth).
func computeMinCost(a []int64, target int64) int64 {
	n := len(a)
	var w [maxDepth]int64
	for d := 0; d < maxDepth; d++ {
		w[d] = (target + (1 << d) - 1) >> d // ceil(target / 2^d)
	}

	var totalW [maxDepth + 1]int64
	for d := 0; d < maxDepth; d++ {
		totalW[d+1] = totalW[d] + w[d]
	}

	best := int64(1<<63 - 1)

	for p := 0; p < n; p++ {
		length := p + 1
		if length > maxDepth {
			length = maxDepth
		}
		sum := int64(0)
		for d := 0; d < length; d++ {
			val := a[p-d]
			if val < w[d] {
				sum += val
			} else {
				sum += w[d]
			}
		}
		cost := totalW[length] - sum
		if cost < best {
			best = cost
		}
	}

	return best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		a := make([]int64, n)
		maxA := int64(0)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] > maxA {
				maxA = a[i]
			}
		}

		qs := make([]int64, q)
		maxK := int64(0)
		for i := 0; i < q; i++ {
			fmt.Fscan(in, &qs[i])
			if qs[i] > maxK {
				maxK = qs[i]
			}
		}

		// Upper bound: each additional pizza increases the target by at most 2 when accounting for the chain,
		// so max reachable value is safely within maxA + 2 * maxK + 5.
		globalHi := maxA + 2*maxK + 5

		low := make([]int64, q)
		high := make([]int64, q)
		for i := 0; i < q; i++ {
			low[i] = maxA
			high[i] = globalHi
		}

		cache := make(map[int64]int64)
		active := true
		for active {
			active = false
			buckets := make(map[int64][]int)
			for i := 0; i < q; i++ {
				if low[i] < high[i] {
					active = true
					mid := (low[i] + high[i] + 1) >> 1
					buckets[mid] = append(buckets[mid], i)
				}
			}
			if !active {
				break
			}

			results := make(map[int64]int64, len(buckets))
			for mid := range buckets {
				if v, ok := cache[mid]; ok {
					results[mid] = v
					continue
				}
				cost := computeMinCost(a, mid)
				cache[mid] = cost
				results[mid] = cost
			}

			for mid, idxs := range buckets {
				cost := results[mid]
				for _, idx := range idxs {
					if cost <= qs[idx] {
						low[idx] = mid
					} else {
						high[idx] = mid - 1
					}
				}
			}
		}

		for i := 0; i < q; i++ {
			fmt.Fprintln(out, low[i])
		}
	}
}
