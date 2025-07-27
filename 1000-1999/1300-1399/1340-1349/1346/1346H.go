package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// This program solves the game described in problemH.txt for contest 1346.
// Players alternately shrink a segment by one unit from the left or right.
// Bob wins immediately after his move if the segment equals one of the
// specified terminal segments. Alice wins if the segment becomes
// degenerate before Bob is able to do so. For an initial segment [l, r],
// Bob can force a win if and only if there exists a terminal segment [L, R]
// such that L + R == l + r and L >= l and R <= r. In that case the number
// of moves Alice performs before defeat equals j = L - l = r - R, which is
// half of the length difference. Therefore we preprocess terminal segments
// grouping them by the sum of their endpoints and, for each sum, store the
// lengths of terminal segments. For an initial segment we select the
// largest terminal length not exceeding its length and sharing the same sum.
// If such a length exists, Bob wins after (len0 - len_t) / 2 moves of Alice;
// otherwise Alice wins (-1).
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	initSeg := make([][2]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &initSeg[i][0], &initSeg[i][1])
	}
	term := make(map[int][]int)
	for i := 0; i < m; i++ {
		var L, R int
		fmt.Fscan(in, &L, &R)
		s := L + R
		l := R - L
		term[s] = append(term[s], l)
	}
	for s := range term {
		sort.Ints(term[s])
	}

	res := make([]int, n)
	for i := 0; i < n; i++ {
		l := initSeg[i][0]
		r := initSeg[i][1]
		sum := l + r
		length := r - l
		lst, ok := term[sum]
		if !ok {
			res[i] = -1
			continue
		}
		// find largest terminal length <= length
		idx := sort.Search(len(lst), func(j int) bool { return lst[j] > length }) - 1
		if idx < 0 {
			res[i] = -1
			continue
		}
		lenT := lst[idx]
		res[i] = (length - lenT) / 2
	}
	for i, v := range res {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
