package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// --- Treap Implementation ---

type TreapNode struct {
	priority    int
	key, val    int
	left, right *TreapNode
}

func NewTreapNode(key, val int) *TreapNode {
	return &TreapNode{
		priority: rand.Int(),
		key:      key,
		val:      val,
	}
}

// split divides the tree into two trees: l (keys <= key) and r (keys > key)
func split(root *TreapNode, key int) (*TreapNode, *TreapNode) {
	if root == nil {
		return nil, nil
	}
	if root.key <= key {
		l, r := split(root.right, key)
		root.right = l
		return root, r
	} else {
		l, r := split(root.left, key)
		root.left = r
		return l, root
	}
}

func merge(l, r *TreapNode) *TreapNode {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	if l.priority > r.priority {
		l.right = merge(l.right, r)
		return l
	} else {
		r.left = merge(l, r.left)
		return r
	}
}

// treapPut updates the value for a key, or inserts it if it doesn't exist.
// 
func treapPut(root *TreapNode, key, val int) *TreapNode {
	// Split into < key, == key, > key
	l, r := split(root, key-1)
	mid, r := split(r, key)

	if mid != nil {
		mid.val = val
	} else {
		mid = NewTreapNode(key, val)
	}

	return merge(merge(l, mid), r)
}

// treapAdd adds w to the value of key, or inserts it if it doesn't exist.
func treapAdd(root *TreapNode, key, w int) *TreapNode {
	l, r := split(root, key-1)
	mid, r := split(r, key)

	if mid != nil {
		mid.val += w
	} else {
		mid = NewTreapNode(key, w)
	}

	return merge(merge(l, mid), r)
}

// treapRemove removes a key from the treap.
func treapRemove(root *TreapNode, key int) *TreapNode {
	l, r := split(root, key-1)
	_, r = split(r, key) // Discard mid
	return merge(l, r)
}

func treapFind(root *TreapNode, key int) int {
	curr := root
	for curr != nil {
		if key == curr.key {
			return curr.val
		}
		if key < curr.key {
			curr = curr.left
		} else {
			curr = curr.right
		}
	}
	return 0
}

func treapContains(root *TreapNode, key int) bool {
	curr := root
	for curr != nil {
		if key == curr.key {
			return true
		}
		if key < curr.key {
			curr = curr.left
		} else {
			curr = curr.right
		}
	}
	return false
}

// treapMaxLessThan returns the node with the largest key strictly less than val
func treapMaxLessThan(root *TreapNode, val int) *TreapNode {
	var best *TreapNode
	curr := root
	for curr != nil {
		if curr.key < val {
			best = curr
			curr = curr.right
		} else {
			curr = curr.left
		}
	}
	return best
}

// treapCeiling returns the node with the smallest key greater than or equal to val
func treapCeiling(root *TreapNode, val int) *TreapNode {
	var best *TreapNode
	curr := root
	for curr != nil {
		if curr.key >= val {
			best = curr
			curr = curr.left
		} else {
			curr = curr.right
		}
	}
	return best
}

func treapTraverse(root *TreapNode, callback func(key, val int)) {
	if root == nil {
		return
	}
	treapTraverse(root.left, callback)
	callback(root.key, root.val)
	treapTraverse(root.right, callback)
}

// --- Priority Queue ---

type Item struct {
	diff int
	keyA int
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].diff != pq[j].diff {
		return pq[i].diff < pq[j].diff
	}
	return pq[i].keyA < pq[j].keyA
}
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(Item))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// --- Main Algorithm Structures ---

type DS struct {
	Q       *PriorityQueue
	A, B    *TreapNode
	tag     int
	cntSize int
}

func NewDS() *DS {
	pq := &PriorityQueue{}
	heap.Init(pq)
	return &DS{
		Q:       pq,
		A:       nil,
		B:       nil,
		tag:     0,
		cntSize: 0,
	}
}

func (ds *DS) Ask(x int) int {
	res := 0
	if treapContains(ds.A, x+ds.tag) {
		res += treapFind(ds.A, x+ds.tag)
	}
	if treapContains(ds.B, x-ds.tag) {
		res += treapFind(ds.B, x-ds.tag)
	}
	return res
}

func (ds *DS) AddA(x int) {
	node := treapMaxLessThan(ds.B, x-2*ds.tag)
	if node != nil {
		heap.Push(ds.Q, Item{diff: x - node.key, keyA: x})
	}
}

func (ds *DS) AddB(x int) {
	node := treapCeiling(ds.A, x+2*ds.tag)
	if node != nil {
		heap.Push(ds.Q, Item{diff: node.key - x, keyA: node.key})
	}
}

func (ds *DS) Add(x, w int) {
	if w > 0 {
		if !treapContains(ds.A, x+ds.tag) {
			ds.cntSize++
			ds.A = treapAdd(ds.A, x+ds.tag, w)
			ds.AddA(x+ds.tag)
		} else {
			ds.A = treapAdd(ds.A, x+ds.tag, w)
		}
	} else {
		if !treapContains(ds.B, x-ds.tag) {
			ds.cntSize++
			ds.B = treapAdd(ds.B, x-ds.tag, w)
			ds.AddB(x-ds.tag)
		} else {
			ds.B = treapAdd(ds.B, x-ds.tag, w)
		}
	}
}

func (ds *DS) Upd(d int) {
	for ds.Q.Len() > 0 {
		top := (*ds.Q)[0]
		if top.diff > 2*(ds.tag+d) {
			break
		}
		heap.Pop(ds.Q)

		ka := top.keyA
		kb := ka - top.diff

		// Verify validity of the event: both keys must exist in their respective maps
		if treapContains(ds.A, ka) && treapContains(ds.B, kb) {
			wa := treapFind(ds.A, ka)
			wb := treapFind(ds.B, kb)
			w := wa + wb

			if w > 0 {
				ds.B = treapRemove(ds.B, kb)
				ds.cntSize--
				ds.A = treapPut(ds.A, ka, w)
				ds.AddA(ka)
			} else {
				ds.A = treapRemove(ds.A, ka)
				ds.cntSize--
				ds.B = treapPut(ds.B, kb, w)
				ds.AddB(kb)
			}
		}
	}
	ds.tag += d
}

// Global Vars
const N = 100005

type kv struct {
	k, v int
}

var (
	ve [N][]struct{ y, z int }
	mp [N]map[int]int
	D  [N]*DS
)

func Dfs(x, fa int) {
	for _, edge := range ve[x] {
		y, z := edge.y, edge.z
		if y == fa {
			continue
		}

		Dfs(y, x)

		type pair struct{ u, v int }
		var now []pair

		for u, v := range mp[y] {
			val1 := -D[y].Ask(u)
			val2 := D[y].Ask(u + 1)
			mx := 0
			if val1 > mx {
				mx = val1
			}
			if val2 > mx {
				mx = val2
			}

			w := v - mx
			if w < 0 {
				w = 0
			}

			if w > 0 {
				now = append(now, pair{u, w})
			}
		}

		D[y].Upd(1)
		for _, p := range now {
			D[y].Add(p.u, p.v)
			D[y].Add(p.u+1, -p.v)
		}
		D[y].Upd(z - 1)

		if D[y].cntSize > D[x].cntSize {
			D[x], D[y] = D[y], D[x]
		}

		treapTraverse(D[y].A, func(u, v int) {
			if v != 0 {
				D[x].Add(u-D[y].tag, v)
			}
		})
		treapTraverse(D[y].B, func(u, v int) {
			if v != 0 {
				D[x].Add(u+D[y].tag, v)
			}
		})
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	rand.Seed(time.Now().UnixNano())

	var n int
	fmt.Fscan(reader, &n)

	for i := 0; i < N; i++ {
		D[i] = NewDS()
		mp[i] = make(map[int]int)
	}

	for i := 1; i < n; i++ {
		var x, y, z int
		fmt.Fscan(reader, &x, &y, &z)
		z *= 2
		ve[x] = append(ve[x], struct{ y, z int }{y, z})
		ve[y] = append(ve[y], struct{ y, z int }{x, z})
	}

	// Add dummy edge 0-1 as per original logic
	ve[0] = append(ve[0], struct{ y, z int }{1, 2})
	ve[1] = append(ve[1], struct{ y, z int }{0, 2})

	var m int
	fmt.Fscan(reader, &m)
	for i := 1; i <= m; i++ {
		var x, y, z int
		fmt.Fscan(reader, &x, &y, &z)
		x *= 2
		if val, exists := mp[z][x]; exists {
			if y > val {
				mp[z][x] = y
			}
		} else {
			mp[z][x] = y
		}
	}

	Dfs(0, 0)

	finalMap := make(map[int]int)
	treapTraverse(D[0].A, func(u, v int) {
		finalMap[u-D[0].tag] += v
	})
	treapTraverse(D[0].B, func(u, v int) {
		finalMap[u+D[0].tag] += v
	})

	var sorted []kv
	for k, v := range finalMap {
		sorted = append(sorted, kv{k, v})
	}

	quickSort(sorted, 0, len(sorted)-1)

	ans := 0
	tmp := 0
	for _, item := range sorted {
		tmp += item.v
		if tmp > ans {
			ans = tmp
		}
	}

	fmt.Fprintln(writer, ans)
}

func quickSort(arr []kv, low, high int) {
	if low < high {
		p := partition(arr, low, high)
		quickSort(arr, low, p-1)
		quickSort(arr, p+1, high)
	}
}

func partition(arr []kv, low, high int) int {
	pivot := arr[high].k
	i := low - 1
	for j := low; j < high; j++ {
		if arr[j].k < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}
