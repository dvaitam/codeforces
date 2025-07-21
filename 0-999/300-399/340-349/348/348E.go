package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "io"
   "os"
)

// Dijkstra implementation
type edge struct { to, w int }
func dijkstra(n, src int, adj [][]edge) []int64 {
   const INF = int64(9e18)
   dist := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = INF
   }
   dist[src] = 0
   // min-heap
   h := &minHeap{{src, 0}}
   heap.Init(h)
   for h.Len() > 0 {
       cur := heap.Pop(h).(pair)
       u, d := cur.u, cur.d
       if d != dist[u] {
           continue
       }
       for _, e := range adj[u] {
           v := e.to
           nd := d + int64(e.w)
           if nd < dist[v] {
               dist[v] = nd
               heap.Push(h, pair{v, nd})
           }
       }
   }
   return dist
}
// heap pair
type pair struct { u int; d int64 }
type minHeap []pair
func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].d < h[j].d }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(pair)) }
func (h *minHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       if err == io.EOF {
           return
       }
       panic(err)
   }
   specials := make([]int, m)
   isSpec := make([]bool, n+1)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &specials[i])
       isSpec[specials[i]] = true
   }
   adj := make([][]edge, n+1)
   for i := 0; i < n-1; i++ {
       var a, b, c int
       fmt.Fscan(in, &a, &b, &c)
       adj[a] = append(adj[a], edge{b, c})
       adj[b] = append(adj[b], edge{a, c})
   }
   // find A
   s0 := specials[0]
   d0 := dijkstra(n, s0, adj)
   A := s0
   for _, u := range specials {
       if d0[u] > d0[A] {
           A = u
       }
   }
   // find B
   dA := dijkstra(n, A, adj)
   B := specials[0]
   for _, u := range specials {
       if dA[u] > dA[B] {
           B = u
       }
   }
   dB := dijkstra(n, B, adj)
   // prepare contribution arrays
   fA := make([]int, n+1)
   cntB := make([]int, n+1)
   // build tree adjacency for rooting
   tree := make([][]int, n+1)
   for u := 1; u <= n; u++ {
       for _, e := range adj[u] {
           tree[u] = append(tree[u], e.to)
       }
   }
   // root tree at A: build parent and depth and binary lifting
   LOG := 17
   parent := make([][]int, LOG)
   for i := 0; i < LOG; i++ {
       parent[i] = make([]int, n+1)
   }
   depth := make([]int, n+1)
   // DFS to set parent[0] and depth
   stack := []struct{u, p int}{{A, 0}}
   for len(stack) > 0 {
       cur := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       u := cur.u; p := cur.p
       parent[0][u] = p
       if p == 0 {
           depth[u] = 0
       } else {
           depth[u] = depth[p] + 1
       }
       for _, v := range tree[u] {
           if v == p {
               continue
           }
           stack = append(stack, struct{u, p int}{v, u})
       }
   }
   // binary lifting
   for k := 1; k < LOG; k++ {
       for v := 1; v <= n; v++ {
           parent[k][v] = parent[k-1][ parent[k-1][v] ]
       }
   }
   // LCA query on tree rooted at A
   lcaA := func(u, v int) int {
       if depth[u] < depth[v] {
           u, v = v, u
       }
       // lift u
       diff := depth[u] - depth[v]
       for k := 0; k < LOG; k++ {
           if diff>>k & 1 == 1 {
               u = parent[k][u]
           }
       }
       if u == v {
           return u
       }
       for k := LOG-1; k >= 0; k-- {
           if parent[k][u] != parent[k][v] {
               u = parent[k][u]
               v = parent[k][v]
           }
       }
       return parent[0][u]
   }
   // classify specials and set contributions
   for _, u := range specials {
       if dA[u] > dB[u] {
           // path to A in root A: all ancestors of u
           fA[u]++
       } else if dB[u] > dA[u] {
           // path to B in root B: count later
           cntB[u]++
       } else {
           // tie: intersection of paths to A and B: ancestors of u up to LCA(u,B)
           L := lcaA(u, B)
           fA[u]++
           // stop after L
           if parent[0][L] != 0 {
               fA[ parent[0][L] ]--
           }
       }
   }
   // DFS sum for root A
   sumA := make([]int, n+1)
   type stt struct{ u, p int; post bool }
   stack := make([]stt, 0, n*2)
   stack = append(stack, stt{A, 0, false})
   for len(stack) > 0 {
       cur := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       if !cur.post {
           // pre
           stack = append(stack, stt{cur.u, cur.p, true})
           for _, v := range tree[cur.u] {
               if v == cur.p {
                   continue
               }
               stack = append(stack, stt{v, cur.u, false})
           }
       } else {
           // post
           s := fA[cur.u]
           for _, v := range tree[cur.u] {
               if v == cur.p {
                   continue
               }
               s += sumA[v]
           }
           sumA[cur.u] = s
       }
   }
   // DFS sum for root B
   sumB := make([]int, n+1)
   stack = stack[:0]
   stack = append(stack, stt{B, 0, false})
   for len(stack) > 0 {
       cur := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       if !cur.post {
           stack = append(stack, stt{cur.u, cur.p, true})
           for _, v := range tree[cur.u] {
               if v == cur.p {
                   continue
               }
               stack = append(stack, stt{v, cur.u, false})
           }
       } else {
           s := cntB[cur.u]
           for _, v := range tree[cur.u] {
               if v == cur.p {
                   continue
               }
               s += sumB[v]
           }
           sumB[cur.u] = s
       }
   }
   // compute answer
   best := 0
   ways := 0
   for v := 1; v <= n; v++ {
       if isSpec[v] {
           continue
       }
       val := sumA[v] + sumB[v]
       if val > best {
           best = val
           ways = 1
       } else if val == best {
           ways++
       }
   }
   fmt.Printf("%d %d\n", best, ways)
}
