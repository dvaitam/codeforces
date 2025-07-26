package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Segment tree for range minimum queries without updates
// Works for arrays of ints

const inf = int(1e9 + 7)

type SegTree struct {
	size int
	tree []int
}

func NewSegTree(arr []int) *SegTree {
	n := len(arr)
	if n == 0 {
		return &SegTree{size: 1, tree: []int{inf, inf}}
	}
	size := 1
	for size < n {
		size <<= 1
	}
	tree := make([]int, 2*size)
	for i := range tree {
		tree[i] = inf
	}
	for i, v := range arr {
		tree[size+i] = v
	}
	for i := size - 1; i > 0; i-- {
		if tree[2*i] < tree[2*i+1] {
			tree[i] = tree[2*i]
		} else {
			tree[i] = tree[2*i+1]
		}
	}
	return &SegTree{size: size, tree: tree}
}

func (s *SegTree) Query(l, r int) int {
	if l > r {
		return inf
	}
	l += s.size
	r += s.size
	res := inf
	for l <= r {
		if l&1 == 1 {
			if s.tree[l] < res {
				res = s.tree[l]
			}
			l++
		}
		if r&1 == 0 {
			if s.tree[r] < res {
				res = s.tree[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func ceilDiv(a int64, b int64) int {
	return int((a + b - 1) / b)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		teachers := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &teachers[i])
		}
		sort.Ints(teachers)
		groups := make([]struct {
			sum  int64
			k    int
			ages []int
			req  int
			pos  int
		}, m)
		totalStudents := 0
		for i := 0; i < m; i++ {
			var k int
			fmt.Fscan(reader, &k)
			groups[i].k = k
			groups[i].ages = make([]int, k)
			totalStudents += k
			var sum int64
			for j := 0; j < k; j++ {
				fmt.Fscan(reader, &groups[i].ages[j])
				sum += int64(groups[i].ages[j])
			}
			groups[i].sum = sum
			groups[i].req = ceilDiv(sum, int64(k))
		}
		// sort groups by required age
		order := make([]int, m)
		for i := 0; i < m; i++ {
			order[i] = i
		}
		sort.Slice(order, func(i, j int) bool {
			if groups[order[i]].req == groups[order[j]].req {
				return order[i] < order[j]
			}
			return groups[order[i]].req < groups[order[j]].req
		})
		reqSorted := make([]int, m)
		for i := 0; i < m; i++ {
			idx := order[i]
			groups[idx].pos = i
			reqSorted[i] = groups[idx].req
		}
		off := n - m
		useT := teachers[off:]
		diff0 := make([]int, m)
		for i := 0; i < m; i++ {
			diff0[i] = useT[i] - reqSorted[i]
		}
		prefMin0 := make([]int, m+1)
		prefMin0[0] = inf
		for i := 0; i < m; i++ {
			if diff0[i] < prefMin0[i] {
				prefMin0[i+1] = diff0[i]
			} else {
				prefMin0[i+1] = prefMin0[i]
			}
		}
		sufMin0 := make([]int, m+1)
		sufMin0[m] = inf
		for i := m - 1; i >= 0; i-- {
			if diff0[i] < sufMin0[i+1] {
				sufMin0[i] = diff0[i]
			} else {
				sufMin0[i] = sufMin0[i+1]
			}
		}
		diff1 := make([]int, 0)
		if m > 1 {
			diff1 = make([]int, m-1)
			for i := 0; i < m-1; i++ {
				diff1[i] = useT[i+1] - reqSorted[i]
			}
		}
		seg1 := NewSegTree(diff1)
		diffNeg1 := make([]int, m)
		diffNeg1[0] = inf
		for i := 1; i < m; i++ {
			diffNeg1[i] = useT[i-1] - reqSorted[i]
		}
		segNeg1 := NewSegTree(diffNeg1)

		res := make([]byte, 0, totalStudents)
		for i := 0; i < m; i++ {
			g := &groups[i]
			for _, age := range g.ages {
				newReq := ceilDiv(g.sum-int64(age), int64(g.k-1))
				idxOld := g.pos
				idxNewPre := sort.Search(len(reqSorted), func(p int) bool { return reqSorted[p] >= newReq })
				if idxNewPre > idxOld {
					idxNewPre--
				}
				idxNew := idxNewPre
				good := true
				if idxNew < idxOld {
					if idxNew > 0 && prefMin0[idxNew] < 0 {
						good = false
					}
					if good && useT[idxNew] < newReq {
						good = false
					}
					if good && idxNew <= idxOld-1 && m > 1 {
						if seg1.Query(idxNew, idxOld-1) < 0 {
							good = false
						}
					}
					if good && idxOld+1 < m {
						if sufMin0[idxOld+1] < 0 {
							good = false
						}
					}
				} else if idxNew > idxOld {
					if idxOld > 0 && prefMin0[idxOld] < 0 {
						good = false
					}
					if good && idxOld+1 <= idxNew && m > 1 {
						if segNeg1.Query(idxOld+1, idxNew) < 0 {
							good = false
						}
					}
					if good && useT[idxNew] < newReq {
						good = false
					}
					if good && idxNew+1 < m {
						if sufMin0[idxNew+1] < 0 {
							good = false
						}
					}
				} else {
					if idxOld > 0 && prefMin0[idxOld] < 0 {
						good = false
					}
					if good && useT[idxOld] < newReq {
						good = false
					}
					if good && idxOld+1 < m {
						if sufMin0[idxOld+1] < 0 {
							good = false
						}
					}
				}
				if good {
					res = append(res, '1')
				} else {
					res = append(res, '0')
				}
			}
		}
		fmt.Fprintln(writer, string(res))
	}
}
