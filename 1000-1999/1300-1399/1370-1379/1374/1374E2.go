package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Item holds time cost and original index
type Item struct {
	time int
	id   int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	v := make([][]Item, 4)
	for i := 1; i <= n; i++ {
		var a, b, c int
		fmt.Fscan(reader, &a, &b, &c)
		state := (b << 0) | (c << 1)
		v[state] = append(v[state], Item{time: a, id: i})
	}
	// sort by time
	for idx := 0; idx < 4; idx++ {
		sort.Slice(v[idx], func(i, j int) bool {
			return v[idx][i].time < v[idx][j].time
		})
	}
	// check feasibility for using c items from v[3]
	check := func(c int) bool {
		if c+(k-c)+(k-c) > m {
			return false
		}
		if k-c > len(v[1]) {
			return false
		}
		if k-c > len(v[2]) {
			return false
		}
		if c+len(v[0])+len(v[1])+len(v[2]) < m {
			return false
		}
		return true
	}
	// initial take from v[3]
	c0 := min(m, len(v[3]))
	if !check(c0) {
		fmt.Fprintln(writer, -1)
		return
	}
	pt := [4]int{}
	var sum int64
	for i := 0; i < c0; i++ {
		sum += int64(v[3][pt[3]].time)
		pt[3]++
	}
	// ensure at least k for both
	for pt[1]+pt[3] < k {
		sum += int64(v[1][pt[1]].time)
		pt[1]++
	}
	for pt[2]+pt[3] < k {
		sum += int64(v[2][pt[2]].time)
		pt[2]++
	}
	// helper to current count
	s := func() int {
		return pt[0] + pt[1] + pt[2] + pt[3]
	}
	// fill to m
	for s() < m {
		best := int64(1e18)
		which := -1
		for i := 0; i <= 2; i++ {
			if pt[i] < len(v[i]) {
				t := int64(v[i][pt[i]].time)
				if t < best {
					best = t
					which = i
				}
			}
		}
		if which == -1 {
			break
		}
		sum += best
		pt[which]++
	}
	ans := sum
	res := pt
	// try reducing v[3] usage
	for pt[3] > 0 {
		pt[3]--
		sum -= int64(v[3][pt[3]].time)
		if !check(pt[3]) {
			break
		}
		for pt[1]+pt[3] < k {
			sum += int64(v[1][pt[1]].time)
			pt[1]++
		}
		for pt[2]+pt[3] < k {
			sum += int64(v[2][pt[2]].time)
			pt[2]++
		}
		// remove excess in v[0]
		for s() > m {
			pt[0]--
			sum -= int64(v[0][pt[0]].time)
		}
		// fill to m again
		for s() < m {
			best := int64(1e18)
			which := -1
			for i := 0; i <= 2; i++ {
				if pt[i] < len(v[i]) {
					t := int64(v[i][pt[i]].time)
					if t < best {
						best = t
						which = i
					}
				}
			}
			if which == -1 {
				break
			}
			sum += best
			pt[which]++
		}
		if sum < ans {
			ans = sum
			res = pt
		}
	}
	// output
	fmt.Fprintln(writer, ans)
	out := make([]int, 0, m)
	for i := 0; i < 4; i++ {
		for j := 0; j < res[i]; j++ {
			out = append(out, v[i][j].id)
		}
	}
	for i, id := range out {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, id)
	}
	writer.WriteByte('\n')
}
