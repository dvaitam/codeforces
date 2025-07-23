package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

type Point struct {
	x int64
	y int64
}

var (
	k, n     int
	stones   []Point
	monsters []Point
	blockers [][][]int
	memo     map[string]bool
)

func between(a, b, c Point) bool {
	dx1 := c.x - a.x
	dy1 := c.y - a.y
	dx2 := b.x - a.x
	dy2 := b.y - a.y
	if dx1*dy2-dy1*dx2 != 0 {
		return false
	}
	if dx1*dx2+dy1*dy2 <= 0 {
		return false
	}
	if dx1*dx1+dy1*dy1 >= dx2*dx2+dy2*dy2 {
		return false
	}
	return true
}

func contains(sl []int, v int) bool {
	for _, x := range sl {
		if x == v {
			return true
		}
	}
	return false
}

func key(list []int, mask uint64) string {
	sort.Ints(list)
	return fmt.Sprintf("%v|%d", list, mask)
}

func killSet(list []int, mask uint64) bool {
	if len(list) == 0 {
		return true
	}
	if bits.OnesCount64(mask) < len(list) {
		return false
	}
	kstr := key(list, mask)
	if val, ok := memo[kstr]; ok {
		return val
	}
	for idx, m := range list {
		rest := append([]int{}, list[:idx]...)
		rest = append(rest, list[idx+1:]...)
		for s := 0; s < k; s++ {
			if mask&(1<<uint(s)) == 0 {
				continue
			}
			bl := blockers[s][m]
			newList := append([]int{}, rest...)
			for _, b := range bl {
				if !contains(newList, b) {
					newList = append(newList, b)
				}
			}
			if len(newList) > bits.OnesCount64(mask)-1 {
				continue
			}
			if killSet(newList, mask^(1<<uint(s))) {
				memo[kstr] = true
				return true
			}
		}
	}
	memo[kstr] = false
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &k, &n); err != nil {
		return
	}
	stones = make([]Point, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &stones[i].x, &stones[i].y)
	}
	monsters = make([]Point, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &monsters[i].x, &monsters[i].y)
	}

	blockers = make([][][]int, k)
	for s := 0; s < k; s++ {
		blockers[s] = make([][]int, n)
		for m := 0; m < n; m++ {
			tmp := make([]int, 0)
			for b := 0; b < n; b++ {
				if b == m {
					continue
				}
				if between(stones[s], monsters[m], monsters[b]) {
					tmp = append(tmp, b)
				}
			}
			blockers[s][m] = tmp
		}
	}

	maskAll := uint64((1 << uint(k)) - 1)
	afraid := 0
	for m := 0; m < n; m++ {
		memo = make(map[string]bool)
		if killSet([]int{m}, maskAll) {
			afraid++
		}
	}
	fmt.Fprintln(writer, afraid)
}
