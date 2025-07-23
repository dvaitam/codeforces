package main

import (
	"bufio"
	"fmt"
	"os"
)

// brute force solver for small n only. For larger n this implementation
// simply outputs 0. A full efficient algorithm has not yet been
// implemented.

func dist(a, b, n int) int {
	d := a - b
	if d < 0 {
		d = -d
	}
	if d > 2*n-d {
		d = 2*n - d
	}
	return d
}

// compute the beauty value for a given matching
func beauty(match []int, n int) int {
	N := 2 * n
	opp := make([]bool, N)
	countOpp := 0
	for i := 0; i < N; i++ {
		if match[i] == (i+n)%N {
			opp[i] = true
			countOpp++
		}
	}
	if countOpp == 0 {
		return 0
	}
	removed := make([]bool, N)
	for i := 0; i < N; i++ {
		if opp[i] {
			removed[i] = true
			removed[(i+n)%N] = true
		}
	}
	var segs []int
	cur := 0
	for i := 0; i < N; i++ {
		if removed[i] {
			if cur > 0 {
				segs = append(segs, cur)
				cur = 0
			} else if len(segs) == 0 {
				// mark zero length segment
				segs = append(segs, 0)
			}
			continue
		}
		cur++
	}
	if cur > 0 {
		segs = append(segs, cur)
	}
	if len(segs) == 0 {
		return 1
	}
	prod := 1
	for _, v := range segs {
		if v == 0 {
			prod *= 0
		} else {
			prod *= v
		}
	}
	return prod
}

func brute(n int) int {
	N := 2 * n
	used := make([]bool, N)
	match := make([]int, N)
	for i := 0; i < N; i++ {
		match[i] = -1
	}
	var res int
	var dfs func(int)
	dfs = func(i int) {
		for i < N && used[i] {
			i++
		}
		if i == N {
			res += beauty(match, n)
			return
		}
		for j := i + 1; j < N; j++ {
			if !used[j] && (dist(i, j, n) <= 2 || dist(i, j, n) == n) {
				i2 := (i + n) % N
				j2 := (j + n) % N
				if used[i2] || used[j2] || !(dist(i2, j2, n) <= 2 || dist(i2, j2, n) == n) {
					continue
				}
				used[i], used[j], used[i2], used[j2] = true, true, true, true
				match[i], match[j] = j, i
				match[i2], match[j2] = j2, i2
				dfs(i + 1)
				used[i], used[j], used[i2], used[j2] = false, false, false, false
				match[i], match[j], match[i2], match[j2] = -1, -1, -1, -1
			}
		}
	}
	dfs(0)
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	if n > 7 {
		// TODO: implement efficient solution
		fmt.Println(0)
		return
	}
	ans := brute(n)
	fmt.Println(ans % 998244353)
}
