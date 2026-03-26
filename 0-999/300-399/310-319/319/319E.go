package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	type Query struct {
		t, x, y int
	}
	queries := make([]Query, n)
	var coords []int

	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &queries[i].t, &queries[i].x, &queries[i].y)
		if queries[i].t == 1 {
			coords = append(coords, queries[i].x, queries[i].y)
		}
	}

	sort.Ints(coords)
	var uniqueCoords []int
	if len(coords) > 0 {
		uniqueCoords = append(uniqueCoords, coords[0])
		for i := 1; i < len(coords); i++ {
			if coords[i] != coords[i-1] {
				uniqueCoords = append(uniqueCoords, coords[i])
			}
		}
	}

	M := len(uniqueCoords)

	getIdx := func(val int) int {
		l, r := 0, M-1
		for l <= r {
			mid := (l + r) / 2
			if uniqueCoords[mid] == val {
				return mid
			}
			if uniqueCoords[mid] < val {
				l = mid + 1
			} else {
				r = mid - 1
			}
		}
		return -1
	}

	tree := make([][]int, 4*M)
	deleted := make([]bool, n+1)
	parent := make([]int, n+1)
	spanL := make([]int, n+1)
	spanR := make([]int, n+1)

	var find func(int) int
	find = func(i int) int {
		if parent[i] == i {
			return i
		}
		parent[i] = find(parent[i])
		return parent[i]
	}

	union := func(a, b int) {
		rootA := find(a)
		rootB := find(b)
		if rootA != rootB {
			parent[rootB] = rootA
		}
	}

	var insertRange func(node, l, r, ql, qr, id int)
	insertRange = func(node, l, r, ql, qr, id int) {
		if ql <= l && r <= qr {
			tree[node] = append(tree[node], id)
			return
		}
		mid := (l + r) / 2
		if ql <= mid {
			insertRange(node*2, l, mid, ql, qr, id)
		}
		if qr > mid {
			insertRange(node*2+1, mid+1, r, ql, qr, id)
		}
	}

	getAll := func(pos int) []int {
		var res []int
		curr := 1
		l, r := 0, M-1
		for {
			if len(tree[curr]) > 0 {
				for _, id := range tree[curr] {
					if !deleted[id] {
						deleted[id] = true
						res = append(res, id)
					}
				}
				tree[curr] = tree[curr][:0]
			}
			if l == r {
				break
			}
			mid := (l + r) / 2
			if pos <= mid {
				curr = curr * 2
				r = mid
			} else {
				curr = curr * 2 + 1
				l = mid + 1
			}
		}
		return res
	}

	intervalCount := 0

	for i := 0; i < n; i++ {
		q := queries[i]
		if q.t == 1 {
			intervalCount++
			id := intervalCount
			parent[id] = id
			L, R := q.x, q.y
			spanL[id] = L
			spanR[id] = R

			idxX := getIdx(q.x)
			idxY := getIdx(q.y)

			compsX := getAll(idxX)
			compsY := getAll(idxY)

			for _, comp := range compsX {
				root := find(comp)
				if root == id {
					continue
				}
				if spanL[root] < L {
					L = spanL[root]
				}
				if spanR[root] > R {
					R = spanR[root]
				}
				union(id, root)
			}
			for _, comp := range compsY {
				root := find(comp)
				if root == id {
					continue
				}
				if spanL[root] < L {
					L = spanL[root]
				}
				if spanR[root] > R {
					R = spanR[root]
				}
				union(id, root)
			}

			root := find(id)
			spanL[root] = L
			spanR[root] = R

			idxL := getIdx(L)
			idxR := getIdx(R)

			if idxL+1 <= idxR-1 {
				insertRange(1, 0, M-1, idxL+1, idxR-1, id)
			}

		} else {
			a, b := q.x, q.y
			rootA := find(a)
			rootB := find(b)
			if rootA == rootB {
				fmt.Fprintln(writer, "YES")
			} else if spanL[rootB] <= spanL[rootA] && spanR[rootA] <= spanR[rootB] {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		}
	}
}
