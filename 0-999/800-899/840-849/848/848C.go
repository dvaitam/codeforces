package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"sort"
)

type IntHeap []int

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type MultiSet struct {
	cnt  map[int]int
	size int
	minH IntHeap
	maxH IntHeap // store negative numbers
}

func NewMultiSet() *MultiSet {
	ms := &MultiSet{cnt: make(map[int]int)}
	heap.Init(&ms.minH)
	heap.Init(&ms.maxH)
	return ms
}

func (ms *MultiSet) Add(x int) {
	heap.Push(&ms.minH, x)
	heap.Push(&ms.maxH, -x)
	ms.cnt[x]++
	ms.size++
}

func (ms *MultiSet) Remove(x int) {
	if ms.cnt[x] > 0 {
		ms.cnt[x]--
		ms.size--
	}
}

func (ms *MultiSet) cleanMin() {
	for ms.minH.Len() > 0 {
		v := ms.minH[0]
		if ms.cnt[v] > 0 {
			break
		}
		heap.Pop(&ms.minH)
	}
}

func (ms *MultiSet) cleanMax() {
	for ms.maxH.Len() > 0 {
		v := -ms.maxH[0]
		if ms.cnt[v] > 0 {
			break
		}
		heap.Pop(&ms.maxH)
	}
}

func (ms *MultiSet) First() int {
	if ms.size == 0 {
		return -1
	}
	ms.cleanMin()
	if ms.minH.Len() == 0 {
		return -1
	}
	return ms.minH[0]
}

func (ms *MultiSet) Last() int {
	if ms.size == 0 {
		return -1
	}
	ms.cleanMax()
	if ms.maxH.Len() == 0 {
		return -1
	}
	return -ms.maxH[0]
}

// We'll implement after.
type Modification struct {
	pos     int
	prevVal int
	newVal  int
}

type Query struct {
	l, r int
	t    int
	idx  int
}

var (
	arr     []int
	sets    map[int]*MultiSet
	curL    int
	curR    int
	curTime int
	mods    []Modification
	memory  int64
)

func getSet(v int) *MultiSet {
	s := sets[v]
	if s == nil {
		s = NewMultiSet()
		sets[v] = s
	}
	return s
}

func addPos(pos int) {
	val := arr[pos]
	s := getSet(val)
	oldContrib := int64(0)
	if s.size > 0 {
		oldContrib = int64(s.Last() - s.First())
	}
	s.Add(pos)
	newContrib := int64(s.Last() - s.First())
	memory += newContrib - oldContrib
}

func removePos(pos int) {
	val := arr[pos]
	s := getSet(val)
	oldContrib := int64(0)
	if s.size > 0 {
		oldContrib = int64(s.Last() - s.First())
	}
	s.Remove(pos)
	newContrib := int64(0)
	if s.size > 0 {
		newContrib = int64(s.Last() - s.First())
	}
	memory += newContrib - oldContrib
}

func applyModification(id int, forward bool) {
	m := mods[id]
	pos := m.pos
	if forward {
		if curL <= pos && pos <= curR {
			removePos(pos)
			arr[pos] = m.newVal
			addPos(pos)
		} else {
			arr[pos] = m.newVal
		}
	} else { // revert
		if curL <= pos && pos <= curR {
			removePos(pos)
			arr[pos] = m.prevVal
			addPos(pos)
		} else {
			arr[pos] = m.prevVal
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	initArr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &initArr[i])
	}
	arr = append([]int(nil), initArr...)
	sets = make(map[int]*MultiSet)
	var queries []Query
	var modCount int
	for i := 0; i < m; i++ {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var p, x int
			fmt.Fscan(in, &p, &x)
			mod := Modification{pos: p, prevVal: arr[p], newVal: x}
			mods = append(mods, mod)
			arr[p] = x
			modCount++
		} else {
			var l, r int
			fmt.Fscan(in, &l, &r)
			queries = append(queries, Query{l: l, r: r, t: modCount, idx: len(queries)})
		}
	}

	// reset array to initial values for processing
	arr = append([]int(nil), initArr...)
	sets = make(map[int]*MultiSet)
	memory = 0
	curL, curR, curTime = 1, 0, 0

	block := int(math.Pow(float64(n), 2.0/3.0))
	if block <= 0 {
		block = 1
	}

	sort.Slice(queries, func(i, j int) bool {
		a, b := queries[i], queries[j]
		blockA := a.l / block
		blockB := b.l / block
		if blockA != blockB {
			return blockA < blockB
		}
		blockRA := a.r / block
		blockRB := b.r / block
		if blockRA != blockRB {
			if blockA%2 == 0 {
				return a.r < b.r
			}
			return a.r > b.r
		}
		return a.t < b.t
	})

	answers := make([]int64, len(queries))
	for _, q := range queries {
		for curTime < q.t {
			applyModification(curTime, true)
			curTime++
		}
		for curTime > q.t {
			curTime--
			applyModification(curTime, false)
		}
		for curL > q.l {
			curL--
			addPos(curL)
		}
		for curR < q.r {
			curR++
			addPos(curR)
		}
		for curL < q.l {
			removePos(curL)
			curL++
		}
		for curR > q.r {
			removePos(curR)
			curR--
		}
		answers[q.idx] = memory
	}

	out := bufio.NewWriter(os.Stdout)
	for i := 0; i < len(answers); i++ {
		fmt.Fprintln(out, answers[i])
	}
	out.Flush()
}
