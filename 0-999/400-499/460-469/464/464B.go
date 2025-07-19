package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

var (
	orig  [8][3]int
	point [8][3]int
	dis   []int64
	miv   []int64
	perms = [6][3]int{
		{0, 1, 2}, {0, 2, 1},
		{1, 0, 2}, {1, 2, 0},
		{2, 0, 1}, {2, 1, 0},
	}
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func dist(i, j int) int64 {
	var d int64
	for k := 0; k < 3; k++ {
		diff := int64(point[i][k] - point[j][k])
		d += diff * diff
	}
	return d
}

func judge() bool {
	if len(dis) != 28 {
		return false
	}
	b := make([]int64, len(dis))
	copy(b, dis)
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	// first 12 equal
	for i := 1; i < 12; i++ {
		if b[i] != b[i-1] {
			return false
		}
	}
	// 13th is double
	if b[12] != 2*b[11] {
		return false
	}
	// next 11 equal
	for i := 13; i < 24; i++ {
		if b[i] != b[i-1] {
			return false
		}
	}
	// 25th is triple of first
	if b[24] != 3*b[0] {
		return false
	}
	// last 3 equal
	for i := 25; i < 28; i++ {
		if b[i] != b[i-1] {
			return false
		}
	}
	return b[0] != 0
}

func dfs(k int) bool {
	if k == 8 {
		return judge()
	}
	// try all permutations of orig[k]
	var backup [3]int
	for pi := 0; pi < 6; pi++ {
		p := perms[pi]
		for j := 0; j < 3; j++ {
			point[k][j] = orig[k][p[j]]
		}
		// record distances to previous points
		ok := true
		// track initial length for backtracking
		before := len(dis)
		for i := 0; i < k; i++ {
			t := dist(i, k)
			// update miv
			if len(miv) == 0 {
				miv = append(miv, t)
			} else {
				miv = append(miv, min(miv[len(miv)-1], t))
			}
			// check prune
			if t > 3*miv[len(miv)-1] {
				ok = false
			}
			dis = append(dis, t)
		}
		if ok && dfs(k+1) {
			return true
		}
		// backtrack
		dis = dis[:before]
		miv = miv[:before]
		// restore not needed since point is overwritten
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < 8; i++ {
		fmt.Fscan(reader, &orig[i][0], &orig[i][1], &orig[i][2])
	}
	if dfs(0) {
		fmt.Println("YES")
		for i := 0; i < 8; i++ {
			fmt.Printf("%d %d %d\n", point[i][0], point[i][1], point[i][2])
		}
	} else {
		fmt.Println("NO")
	}
}
