package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod2164F1 = 998244353
const combLimit2164F1 = 6005

var fact2164F1 [combLimit2164F1]int
var invFact2164F1 [combLimit2164F1]int

func init() {
	fact2164F1[0] = 1
	for i := 1; i < combLimit2164F1; i++ {
		fact2164F1[i] = int(int64(fact2164F1[i-1]) * int64(i) % mod2164F1)
	}
	invFact2164F1[combLimit2164F1-1] = pow2164F1(fact2164F1[combLimit2164F1-1], mod2164F1-2)
	for i := combLimit2164F1 - 2; i >= 0; i-- {
		invFact2164F1[i] = int(int64(invFact2164F1[i+1]) * int64(i+1) % mod2164F1)
	}
}

func pow2164F1(a, e int) int {
	res := 1
	base := a % mod2164F1
	exp := e
	for exp > 0 {
		if exp&1 == 1 {
			res = int(int64(res) * int64(base) % mod2164F1)
		}
		base = int(int64(base) * int64(base) % mod2164F1)
		exp >>= 1
	}
	return res
}

func comb2164F1(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	return int(int64(fact2164F1[n]) * int64(invFact2164F1[k]) % mod2164F1 * int64(invFact2164F1[n-k]) % mod2164F1)
}

func mul2164F1(a, b int) int {
	return int(int64(a) * int64(b) % mod2164F1)
}

type Treap2164F1 struct {
	left, right, size, prio, val []int
	root, tot                    int
	seed                         uint64
}

func newTreap2164F1(capacity int) *Treap2164F1 {
	return &Treap2164F1{
		left:  make([]int, capacity),
		right: make([]int, capacity),
		size:  make([]int, capacity),
		prio:  make([]int, capacity),
		val:   make([]int, capacity),
		seed:  1,
	}
}

func (t *Treap2164F1) nextRand() int {
	t.seed ^= t.seed << 7
	t.seed ^= t.seed >> 9
	return int(t.seed & 0x7fffffff)
}

func (t *Treap2164F1) update(x int) {
	if x == 0 {
		return
	}
	t.size[x] = 1 + t.size[t.left[x]] + t.size[t.right[x]]
}

func (t *Treap2164F1) split(x, k int) (int, int) {
	if x == 0 {
		return 0, 0
	}
	leftSize := t.size[t.left[x]]
	if k > leftSize {
		leftSub, rightSub := t.split(t.right[x], k-leftSize-1)
		t.right[x] = leftSub
		t.update(x)
		return x, rightSub
	}
	leftSub, rightSub := t.split(t.left[x], k)
	t.left[x] = rightSub
	t.update(x)
	return leftSub, x
}

func (t *Treap2164F1) merge(a, b int) int {
	if a == 0 || b == 0 {
		return a + b
	}
	if t.prio[a] < t.prio[b] {
		t.right[a] = t.merge(t.right[a], b)
		t.update(a)
		return a
	}
	t.left[b] = t.merge(a, t.left[b])
	t.update(b)
	return b
}

func (t *Treap2164F1) newNode(val int) int {
	t.tot++
	idx := t.tot
	t.val[idx] = val
	t.prio[idx] = t.nextRand()
	t.left[idx], t.right[idx] = 0, 0
	t.size[idx] = 1
	return idx
}

func (t *Treap2164F1) insert(val, pos int) {
	left, right := t.split(t.root, pos)
	node := t.newNode(val)
	t.root = t.merge(t.merge(left, node), right)
}

func (t *Treap2164F1) delete(pos int) {
	left, right := t.split(t.root, pos)
	_, right = t.split(right, 1)
	t.root = t.merge(left, right)
}

func (t *Treap2164F1) kth(k int) int {
	x := t.root
	for x != 0 {
		leftSize := t.size[t.left[x]]
		if k <= leftSize {
			x = t.left[x]
		} else if k == leftSize+1 {
			return t.val[x]
		} else {
			k -= leftSize + 1
			x = t.right[x]
		}
	}
	return -1
}

func pairKey2164F1(u, v int) int64 {
	return (int64(u) << 32) | int64(uint32(v))
}

func solve2164F1(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	fmt.Fscan(reader, &n)
	origN := n
	sentinel := origN + 1
	adj := make([][]int, sentinel+1)
	for i := 2; i <= origN; i++ {
		var p int
		fmt.Fscan(reader, &p)
		adj[p] = append(adj[p], i)
	}
	a := make([]int, sentinel+1)
	for i := 1; i <= origN; i++ {
		fmt.Fscan(reader, &a[i])
		a[i]++
	}
	dep := make([]int, sentinel+1)
	dep[1] = 1
	var dfsDepth func(int)
	dfsDepth = func(u int) {
		for _, v := range adj[u] {
			dep[v] = dep[u] + 1
			dfsDepth(v)
		}
	}
	dfsDepth(1)
	for i := 1; i <= origN; i++ {
		if a[i] > dep[i] {
			fmt.Fprintln(writer, 0)
			return
		}
	}

	treap := newTreap2164F1(origN + 5)
	treap.insert(0, 0)
	treap.insert(sentinel, 1)

	S := make([]int, origN+1)
	T := make([]int, origN+1)
	dfn := make([]int, origN+1)
	leftID := make([]int, origN+1)
	rightID := make([]int, origN+1)
	seq := make([]int, 0, origN)
	edgeID := make(map[int64]int, 2*origN+5)
	edgeID[pairKey2164F1(0, sentinel)] = 1

	maxEdges := 2*origN + 5
	F := make([]int, maxEdges)
	G := make([]int, maxEdges)
	F[1] = 1

	idx := 0
	var dfs func(int)
	dfs = func(u int) {
		idx++
		dfn[u] = idx
		s := treap.kth(a[u])
		t := treap.kth(a[u] + 1)
		S[u] = s
		T[u] = t
		left := dfn[u] << 1
		right := left | 1
		leftID[u] = left
		rightID[u] = right
		edgeID[pairKey2164F1(s, u)] = left
		edgeID[pairKey2164F1(u, t)] = right
		F[left] = 1
		F[right] = 1
		treap.insert(u, a[u])
		seq = append(seq, u)
		for _, v := range adj[u] {
			dfs(v)
		}
		treap.delete(a[u])
	}
	dfs(1)

	for i := len(seq) - 1; i >= 0; i-- {
		x := seq[i]
		s := S[x]
		t := T[x]
		id := edgeID[pairKey2164F1(s, t)]
		left := leftID[x]
		right := rightID[x]
		fst := F[id]
		gst := G[id]
		fsx := F[left]
		gsx := G[left]
		fxt := F[right]
		gxt := G[right]
		total := gst + gsx + gxt + 1
		combVal := comb2164F1(total, gst)
		val := mul2164F1(mul2164F1(mul2164F1(fst, fsx), fxt), combVal)
		F[id] = val
		G[id] = total
	}

	fmt.Fprintln(writer, F[edgeID[pairKey2164F1(0, sentinel)]])
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solve2164F1(reader, writer)
	}
	writer.Flush()
}
