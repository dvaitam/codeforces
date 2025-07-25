package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Pair maps a node to a centroid component
type Pair struct{ cid, idx int }

// HeapNode holds a lazy heap element
type HeapNode struct{ val int64; idx int }

// MaxHeap implements a max-heap of HeapNode
type MaxHeap []HeapNode
func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].val > h[j].val }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(HeapNode)) }
func (h *MaxHeap) Pop() interface{} {
   old := *h; n := len(old)
   x := old[n-1]; *h = old[:n-1]
   return x
}

// SegTree for range max query and point update
type SegTree struct{ n int; t []int64 }
func NewSegTree(a []int64) *SegTree {
   n := len(a); t := make([]int64, 2*n)
   for i := 0; i < n; i++ { t[n+i] = a[i] }
   for i := n - 1; i > 0; i-- { if t[i<<1] > t[i<<1|1] { t[i] = t[i<<1] } else { t[i] = t[i<<1|1] } }
   return &SegTree{n, t}
}
func (s *SegTree) Update(p int, v int64) {
   p += s.n; s.t[p] = v
   for p >>= 1; p > 0; p >>= 1 {
       if s.t[p<<1] > s.t[p<<1|1] { s.t[p] = s.t[p<<1] } else { s.t[p] = s.t[p<<1|1] }
   }
}
func (s *SegTree) Query() int64 { return s.t[1] }

var (
   n, q int
   Wlimit int64
   adj [][]Edge
   edges []EdgeW
   deleted []bool
   sz []int
   centroidID []int
   nextCID int
   // mapping for each node to its centroid components
   vpairs [][]Pair
   // per centroid component distances
   distComp [][]int64
   heaps []*MaxHeap
)

// Edge in adj list
type Edge struct{ to, id int }
// EdgeW stores endpoints and weight
type EdgeW struct{ u, v int; w int64 }

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n, &q, &Wlimit)
   adj = make([][]Edge, n+1)
   edges = make([]EdgeW, n-1)
   for i := 0; i < n-1; i++ {
       var u, v int; var w int64
       fmt.Fscan(in, &u, &v, &w)
       edges[i] = EdgeW{u, v, w}
       adj[u] = append(adj[u], Edge{v, i})
       adj[v] = append(adj[v], Edge{u, i})
   }
   // parent and child mapping
   parent := make([]int, n+1)
   child := make([]int, n-1)
   // stack DFS
   st := []int{1}
   parent[1] = -1
   for len(st) > 0 {
       u := st[len(st)-1]; st = st[:len(st)-1]
       for _, e := range adj[u] {
           v := e.to
           if v == parent[u] { continue }
           parent[v] = u
           child[e.id] = v
           st = append(st, v)
       }
   }
   // init centroid globals
   deleted = make([]bool, n+1)
   sz = make([]int, n+1)
   centroidID = make([]int, n+1)
   vpairs = make([][]Pair, n+1)
   distComp = [][]int64{}
   heaps = []*MaxHeap{}
   nextCID = 0
   // decompose
   decompose(1)
   // init heaps and initial dia array
   dia := make([]int64, nextCID)
   for cid := 0; cid < nextCID; cid++ {
       h := &MaxHeap{}
       heap.Init(h)
       for idx, d := range distComp[cid] {
           heap.Push(h, HeapNode{d, idx})
       }
       heaps = append(heaps, h)
       // get top two
       m1 := popTop(h)
       m2 := popTop(h)
       dia[cid] = m1 + m2
       // push back
       heap.Push(h, HeapNode{m1, 0})
       heap.Push(h, HeapNode{m2, 1})
   }
   // segment tree
   // build size as power of two
   size := 1
   for size < nextCID { size <<= 1 }
   // pad dia
   arr := make([]int64, size)
   for i := 0; i < nextCID; i++ { arr[i] = dia[i] }
   stree := NewSegTree(arr)
   // process queries
   last := int64(0)
   for i := 0; i < q; i++ {
       var d int; var e int64
       fmt.Fscan(in, &d, &e)
       d = int((int64(d)+last) % int64(n-1))
       e = (e + last) % Wlimit
       delta := e - edges[d].w
       edges[d].w = e
       if delta != 0 {
           v := child[d]
           for _, p := range vpairs[v] {
               cid, idx := p.cid, p.idx
               distComp[cid][idx] += delta
               heap.Push(heaps[cid], HeapNode{distComp[cid][idx], idx})
               // recompute dia for cid
               m1 := cleanPop(heaps[cid], distComp[cid])
               m2 := cleanPop(heaps[cid], distComp[cid])
               newd := m1 + m2
               stree.Update(cid, newd)
               // push back tops
               heap.Push(heaps[cid], HeapNode{m1, 0})
               heap.Push(heaps[cid], HeapNode{m2, 0})
           }
       }
       ans := stree.Query()
       fmt.Fprintln(out, ans)
       last = ans
   }
}

// pop top or return 0
func popTop(h *MaxHeap) int64 {
   if h.Len() == 0 { return 0 }
   return heap.Pop(h).(HeapNode).val
}

// cleanPop pops invalid until matches current distComp and returns val
func cleanPop(h *MaxHeap, dist []int64) int64 {
   for h.Len() > 0 {
       top := (*h)[0]
       if dist[top.idx] == top.val {
           heap.Pop(h)
           return top.val
       }
       heap.Pop(h)
   }
   return 0
}

// dfsSize computes subtree sizes
func dfsSize(u, p int) {
   sz[u] = 1
   for _, e := range adj[u] {
       v := e.to
       if v == p || deleted[v] { continue }
       dfsSize(v, u)
       sz[u] += sz[v]
   }
}

// dfsCentroid finds centroid
func dfsCentroid(u, p, total int) int {
   for _, e := range adj[u] {
       v := e.to
       if v == p || deleted[v] { continue }
       if sz[v] > total/2 {
           return dfsCentroid(v, u, total)
       }
   }
   return u
}

// dfsCollect collects nodes in component
func dfsCollect(u, p, cid, idx int, dist int64) {
   // record mapping
   vpairs[u] = append(vpairs[u], Pair{cid, idx})
   if dist > distComp[cid][idx] {
       distComp[cid][idx] = dist
   }
   for _, e := range adj[u] {
       v := e.to
       if v == p || deleted[v] { continue }
       w := edges[e.id].w
       dfsCollect(v, u, cid, idx, dist + w)
   }
}

// decompose centroid
func decompose(u int) {
   dfsSize(u, -1)
   c := dfsCentroid(u, -1, sz[u])
   cid := nextCID
   centroidID[c] = cid
   nextCID++
   // prepare distComp for cid
   // count comps by iterating neighbors
   cnt := 0
   for _, e := range adj[c] {
       if deleted[e.to] { continue }
       cnt++
   }
   distComp = append(distComp, make([]int64, cnt))
   // collect for each component
   idx := 0
   deleted[c] = true
   for _, e := range adj[c] {
       v := e.to
       if deleted[v] { continue }
       w := edges[e.id].w
       dfsCollect(v, c, cid, idx, w)
       idx++
   }
   // decompose subtrees
   for _, e := range adj[c] {
       v := e.to
       if deleted[v] { continue }
       decompose(v)
   }
}
