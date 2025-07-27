package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Dynamic tree 1-center: maintain dynamic diameter of centers of paths
// we subdivide each edge and build centroid decomposition

type CDPair struct{ c, dist int }

// Max-heap of ints
type IntHeap []int
func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() interface{} {
   old := *h; n := len(old); x := old[n-1]; *h = old[0 : n-1]; return x
}

// Multiset with lazy deletion
type MultiSet struct{ inc, del *IntHeap }
func NewMultiSet() *MultiSet {
   inc, del := &IntHeap{}, &IntHeap{}
   heap.Init(inc); heap.Init(del)
   return &MultiSet{inc, del}
}
func (ms *MultiSet) Insert(x int) { heap.Push(ms.inc, x) }
func (ms *MultiSet) Remove(x int) { heap.Push(ms.del, x) }
func (ms *MultiSet) clean() {
   for ms.del.Len() > 0 && ms.inc.Len() > 0 && (*ms.inc)[0] == (*ms.del)[0] {
       heap.Pop(ms.inc); heap.Pop(ms.del)
   }
}
func (ms *MultiSet) Max() int {
   ms.clean()
   if ms.inc.Len() == 0 { return 0 }
   return (*ms.inc)[0]
}
func (ms *MultiSet) SecondMax() int {
   ms.clean()
   if ms.inc.Len() < 2 { return 0 }
   // pop first, get second, push back first
   first := heap.Pop(ms.inc).(int)
   ms.clean()
   second := 0
   if ms.inc.Len() > 0 { second = (*ms.inc)[0] }
   heap.Push(ms.inc, first)
   return second
}

// Global heap for centroid sums
type CHItem struct{ sum, c int }
type CH []*CHItem
func (h CH) Len() int           { return len(h) }
func (h CH) Less(i, j int) bool { return h[i].sum > h[j].sum }
func (h CH) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *CH) Push(x interface{}) { *h = append(*h, x.(*CHItem)) }
func (h *CH) Pop() interface{} {
   old := *h; n := len(old); x := old[n-1]; *h = old[0 : n-1]; return x
}

var (
   origAdj    [][]struct{ to, mid int }
   N, Q       int
   biParent   [][]int
   depthOrig  []int
   lmax       int
   adj        [][]int
   N2         int
   removed    []bool
   sz         []int
   parentCD   []int
   cdAnc      [][]CDPair
   msList     []*MultiSet
   curSum     []int
   globalCH   = &CH{}
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &N, &Q)
   origAdj = make([][]struct{ to, mid int }, N+1)
   edges := make([][2]int, N-1)
   for i := 0; i < N-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       edges[i] = [2]int{u, v}
   }
   // build subdivided tree
   N2 = N + (N - 1)
   adj = make([][]int, N2+1)
   nextId := N + 1
   for _, e := range edges {
       u, v := e[0], e[1]
       w := nextId; nextId++
       adj[u] = append(adj[u], w)
       adj[w] = append(adj[w], u)
       adj[v] = append(adj[v], w)
       adj[w] = append(adj[w], v)
       origAdj[u] = append(origAdj[u], struct{to, mid int}{v, w})
       origAdj[v] = append(origAdj[v], struct{to, mid int}{u, w})
   }
   // LCA on original
   lmax = 1
   for (1<<lmax) <= N { lmax++ }
   biParent = make([][]int, lmax)
   for i := range biParent {
       biParent[i] = make([]int, N+1)
   }
   depthOrig = make([]int, N+1)
   // build original tree for LCA
   treeOrig := make([][]int, N+1)
   for _, e := range edges {
       treeOrig[e[0]] = append(treeOrig[e[0]], e[1])
       treeOrig[e[1]] = append(treeOrig[e[1]], e[0])
   }
   var dfsLCA func(u, p int)
   dfsLCA = func(u, p int) {
       biParent[0][u] = p
       for _, v := range treeOrig[u] {
           if v == p { continue }
           depthOrig[v] = depthOrig[u] + 1
           dfsLCA(v, u)
       }
   }
   dfsLCA(1, 0)
   for i := 1; i < lmax; i++ {
       for v := 1; v <= N; v++ {
           biParent[i][v] = biParent[i-1][ biParent[i-1][v] ]
       }
   }
   // centroid decomposition init
   removed = make([]bool, N2+1)
   sz = make([]int, N2+1)
   parentCD = make([]int, N2+1)
   cdAnc = make([][]CDPair, N2+1)
   buildCD(1, 0)
   // init multisets and global heap
   msList = make([]*MultiSet, N2+1)
   curSum = make([]int, N2+1)
   for i := 1; i <= N2; i++ {
       msList[i] = NewMultiSet()
   }
   heap.Init(globalCH)
   // process queries
   for i := 0; i < Q; i++ {
       var tp int
       fmt.Fscan(in, &tp)
       if tp == 1 || tp == 2 {
           var u, v int
           fmt.Fscan(in, &u, &v)
           // find center node in subdivided
           duv := depthOrig[u] + depthOrig[v] - 2*depthOrig[lca(u, v)]
           // duv orig distance
           k := duv / 2
           mid := kthNode(u, v, k)
           center := mid
           if duv%2 == 1 {
               // next node
               next := kthNode(u, v, k+1)
               // find subdiv node
               for _, e := range origAdj[mid] {
                   if e.to == next {
                       center = e.mid
                       break
                   }
               }
           }
           if tp == 1 {
               activate(center, 1)
           } else {
               activate(center, -1)
           }
       } else {
           var d int
           fmt.Fscan(in, &d)
           // get current diameter doubled
           var diam int
           for globalCH.Len() > 0 {
               top := (*globalCH)[0]
               if curSum[top.c] != top.sum {
                   heap.Pop(globalCH)
                   continue
               }
               diam = top.sum
               break
           }
           if diam <= 2*d {
               out.WriteString("Yes\n")
           } else {
               out.WriteString("No\n")
           }
       }
   }
}

func lca(u, v int) int {
   if depthOrig[u] < depthOrig[v] { u, v = v, u }
   diff := depthOrig[u] - depthOrig[v]
   for i := 0; diff > 0; i++ {
       if diff&1 == 1 { u = biParent[i][u] }
       diff >>= 1
   }
   if u == v { return u }
   for i := lmax - 1; i >= 0; i-- {
       if biParent[i][u] != biParent[i][v] {
           u = biParent[i][u]
           v = biParent[i][v]
       }
   }
   return biParent[0][u]
}

// move k steps from u towards v
func kthNode(u, v, k int) int {
   w := lca(u, v)
   duw := depthOrig[u] - depthOrig[w]
   if k <= duw {
       return ancestor(u, k)
   }
   k2 := depthOrig[v] - depthOrig[w] - (k - duw)
   return ancestor(v, k2)
}

func ancestor(u, k int) int {
   for i := 0; k > 0; i++ {
       if k&1 == 1 { u = biParent[i][u] }
       k >>= 1
   }
   return u
}

func activate(x, delta int) {
   for _, p := range cdAnc[x] {
       c, d := p.c, p.dist
       if delta > 0 {
           msList[c].Insert(d)
       } else {
           msList[c].Remove(d)
       }
       sum := msList[c].Max() + msList[c].SecondMax()
       curSum[c] = sum
       heap.Push(globalCH, &CHItem{sum, c})
   }
}

// centroid decomposition building
func buildCD(start, p int) {
   nodes := collect(start, -1)
   c := findCentroid(nodes, -1, len(nodes))
   parentCD[c] = p
   // record distances from c
   dfsCD(c, -1, 0, c)
   removed[c] = true
   for _, v := range adj[c] {
       if !removed[v] {
           buildCD(v, c)
       }
   }
}

func collect(u, p int) []int {
   var list []int
   var dfs func(int, int)
   dfs = func(u, p int) {
       list = append(list, u)
       for _, v := range adj[u] {
           if v == p || removed[v] { continue }
           dfs(v, u)
       }
   }
   dfs(u, p)
   return list
}

func findCentroid(nodes []int, p, total int) int {
   var dfsSize func(int, int) int
   dfsSize = func(u, p int) int {
       sz[u] = 1
       for _, v := range adj[u] {
           if v == p || removed[v] { continue }
           sz[u] += dfsSize(v, u)
       }
       return sz[u]
   }
   dfsSize(nodes[0], -1)
   u := nodes[0]
   for {
       heavy := -1
       for _, v := range adj[u] {
           if removed[v] || sz[v] > sz[u] || sz[v] > total/2 { continue }
           if heavy < 0 || sz[v] > sz[heavy] { heavy = v }
       }
       if heavy < 0 { break }
       u = heavy
   }
   return u
}

func dfsCD(u, p, d, cen int) {
   cdAnc[u] = append(cdAnc[u], CDPair{cen, d})
   for _, v := range adj[u] {
       if v == p || removed[v] { continue }
       dfsCD(v, u, d+1, cen)
   }
}
