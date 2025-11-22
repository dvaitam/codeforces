package main

import (
	"bufio"
	"fmt"
	"os"
)

// The graph is a simple cycle with an extra edge (x, y); only x and y have
// degree 3, all other nodes have degree 2.  For a degree-2 node with neighbours
// u and v, its value must be mex(value[u], value[v]), which restricts values to
// {0,1,2}.  For x and y (degree 3), the value can also be 3.  Overall all values
// stay within 0..3.
//
// We brute force values for x and y (16 combinations) and ask whether both
// paths between them can be filled with degree-2 constraints, and whether the
// degree-3 constraints at x and y hold.  The degree-2 path feasibility is
// checked with a small-state DP over pairs of consecutive values; domain is
// only 0..3 so transitions are constant time.
//
// Total time is O(n * 16) across a test case, easily within limits.

type state struct {
	ok bool
	p  int8
}

// nextOptions returns all values t such that mex(prev, t) == cur.
func nextOptions(prev, cur int) (arr [4]int, cnt int) {
	switch cur {
	case 0:
		if prev == 0 {
			return
		}
		arr = [4]int{1, 2, 3}
		return arr, 3
	case 1:
		if prev == 1 {
			return
		}
		// need 0 present, 1 absent
		if prev == 0 {
			arr = [4]int{0, 2, 3}
			return arr, 3
		}
		// prev in {2,3}, so next must be 0
		arr[0] = 0
		return arr, 1
	case 2:
		if prev == 0 {
			arr[0] = 1
			return arr, 1
		}
		if prev == 1 {
			arr[0] = 0
			return arr, 1
		}
		return
	default: // cur == 3 impossible for degree-2 constraints
		return
	}
}

func mex3(a, b, c int) int {
	seen := [4]bool{}
	seen[a], seen[b], seen[c] = true, true, true
	for i := 0; i < 4; i++ {
		if !seen[i] {
			return i
		}
	}
	return 4
}

// buildPath checks if a path of len edges (len+1 nodes) can be filled between
// endpoints start and end. On success returns the node values in order from
// start to end.
func buildPath(length int, start, end int) (bool, []int) {
	L := length
	values := make([]int, L+1)

	if L == 1 {
		values[0], values[1] = start, end
		return true, values
	}

	var cur, nxt [4][4]bool
	parents := make([][4][4]int8, L+1)
	for i := 0; i <= L; i++ {
		for p := 0; p < 4; p++ {
			for c := 0; c < 4; c++ {
				parents[i][p][c] = -1
			}
		}
	}

	for v := 0; v <= 2; v++ { // node 1 is degree 2 unless L==1 handled above
		cur[start][v] = true
	}

	for i := 1; i <= L-1; i++ {
		for p := 0; p < 4; p++ {
			for c := 0; c < 4; c++ {
				nxt[p][c] = false
			}
		}
		lastIdx := i + 1
		for p := 0; p < 4; p++ {
			for c := 0; c < 4; c++ {
				if !cur[p][c] {
					continue
				}
				opts, cnt := nextOptions(p, c)
				for k := 0; k < cnt; k++ {
					nv := opts[k]
					if lastIdx != L && nv >= 3 {
						continue // internal node cannot be 3
					}
					if lastIdx == L && nv != end {
						continue
					}
					if !nxt[c][nv] {
						nxt[c][nv] = true
						parents[lastIdx][c][nv] = int8(p)
					}
				}
			}
		}
		cur = nxt
	}

	prev := -1
	for p := 0; p < 4; p++ {
		if cur[p][end] {
			prev = p
			break
		}
	}
	if prev == -1 {
		return false, nil
	}

	values[L] = end
	values[L-1] = prev
	for idx := L; idx >= 2; idx-- {
		p := parents[idx][values[idx-1]][values[idx]]
		if p < 0 {
			return false, nil
		}
		values[idx-2] = int(p)
	}
	return true, values
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
		var x, y int
		fmt.Fscan(in, &n, &x, &y)
		x--
		y--
		if x == y {
			// not possible by constraints, but guard anyway
			for i := 0; i < n; i++ {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, 0)
			}
			fmt.Fprintln(out)
			continue
		}

		d := (y - x + n) % n
		len1 := d     // edges clockwise from x to y
		len2 := n - d // edges counter-clockwise
		idx1 := make([]int, len1+1)
		for i := 0; i <= len1; i++ {
			idx1[i] = (x + i) % n
		}
		idx2 := make([]int, len2+1)
		for i := 0; i <= len2; i++ {
			idx2[i] = (x - i%n + n) % n
		}

		found := false
		var ans []int

		for vx := 0; vx <= 3 && !found; vx++ {
			for vy := 0; vy <= 3 && !found; vy++ {
				ok1, path1 := buildPath(len1, vx, vy)
				if !ok1 {
					continue
				}
				ok2, path2 := buildPath(len2, vx, vy)
				if !ok2 {
					continue
				}

				// neighbours for x
				nx1 := path1[1]
				nx2 := path2[1]
				if mex3(nx1, nx2, vy) != vx {
					continue
				}
				// neighbours for y
				ny1 := path1[len1-1]
				ny2 := path2[len2-1]
				if mex3(ny1, ny2, vx) != vy {
					continue
				}

				ans = make([]int, n)
				for i, idx := range idx1 {
					ans[idx] = path1[i]
				}
				for i := 1; i < len2; i++ { // internal nodes on second path
					ans[idx2[i]] = path2[i]
				}
				found = true
			}
		}

		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
