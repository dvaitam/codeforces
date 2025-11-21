package main

import (
	"bufio"
	"fmt"
	"os"
)

type solver struct {
	cnt   map[int]int
	edges map[int][]int
	memo  map[int]int
}

func (s *solver) dp(x int) int {
	if v, ok := s.memo[x]; ok {
		return v
	}
	res := x
	for _, v := range s.edges[x] {
		curr := v
		if s.cnt[v] > 0 {
			curr = s.dp(v)
		}
		if curr > res {
			res = curr
		}
	}
	s.memo[x] = res
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		var m int
		fmt.Fscan(in, &n, &m)

		cnt := make(map[int]int)
		edges := make(map[int][]int)
		maxMex := 0

		for i := 0; i < n; i++ {
			var l int
			fmt.Fscan(in, &l)
			values := make(map[int]struct{}, l)
			for j := 0; j < l; j++ {
				var x int
				fmt.Fscan(in, &x)
				values[x] = struct{}{}
			}

			mex := 0
			for {
				if _, ok := values[mex]; ok {
					mex++
				} else {
					break
				}
			}
			mex2 := mex + 1
			for {
				if _, ok := values[mex2]; ok {
					mex2++
				} else {
					break
				}
			}

			cnt[mex]++
			edges[mex] = append(edges[mex], mex2)
			if mex > maxMex {
				maxMex = mex
			}
		}

		s := solver{
			cnt:   cnt,
			edges: edges,
			memo:  make(map[int]int),
		}

		// Precompute dp for all mex values that appear.
		for x := range cnt {
			s.dp(x)
		}

		// generalMax: maximum reachable using sequences with mex count >= 2 (they can activate themselves).
		generalMax := maxMex
		seen := make(map[int]bool)
		queue := make([]int, 0)
		for x, c := range cnt {
			if c >= 2 {
				queue = append(queue, x)
			}
		}
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			if seen[u] {
				continue
			}
			seen[u] = true
			if val := s.dp(u); val > generalMax {
				generalMax = val
			}
			for _, v := range edges[u] {
				if cnt[v] > 0 && !seen[v] {
					queue = append(queue, v)
				}
			}
		}

		limit := maxMex
		if m < limit {
			limit = m
		}

		var ansSum int64
		for k := 0; k <= limit; k++ {
			base := maxMex
			if k > base {
				base = k
			}
			res := generalMax
			if base > res {
				res = base
			}
			if cnt[k] > 0 {
				if val := s.dp(k); val > res {
					res = val
				}
			}
			ansSum += int64(res)
		}

		if m > maxMex {
			low := maxMex + 1
			high := m
			if generalMax > maxMex {
				mid := generalMax - 1
				if mid > high {
					mid = high
				}
				if mid >= low {
					count := int64(mid-low+1)
					ansSum += count * int64(generalMax)
					low = mid + 1
				}
			}
			if low <= high {
				a := int64(low)
				b := int64(high)
				count := b - a + 1
				ansSum += (a + b) * count / 2
			}
		}

		fmt.Fprintln(out, ansSum)
	}
}
