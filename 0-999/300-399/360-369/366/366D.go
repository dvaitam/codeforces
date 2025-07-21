package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

type Edge struct {
   to, rk int
}

// Item for max-heap: node with current width
type Item struct {
   node, w int
}

type MaxHeap []Item

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].w > h[j].w }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   item := old[n-1]
   *h = old[0 : n-1]
   return item
}

func widestPath(n int, adj [][]Edge) int {
   const INF = 1000000001
   width := make([]int, n+1)
   for i := 1; i <= n; i++ {
       width[i] = 0
   }
   width[1] = INF
   h := &MaxHeap{}
   heap.Init(h)
   heap.Push(h, Item{node: 1, w: INF})
   for h.Len() > 0 {
       it := heap.Pop(h).(Item)
       u, w := it.node, it.w
       if w < width[u] {
           continue
       }
       if u == n {
           break
       }
       for _, e := range adj[u] {
           nw := w
           if e.rk < nw {
               nw = e.rk
           }
           if nw > width[e.to] {
               width[e.to] = nw
               heap.Push(h, Item{node: e.to, w: nw})
           }
       }
   }
   return width[n]
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   type E struct{ u, v, l, r int }
   edges := make([]E, m)
   ls := make([]int, 0, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &edges[i].u, &edges[i].v, &edges[i].l, &edges[i].r)
       ls = append(ls, edges[i].l)
   }
   if m == 0 {
       fmt.Println("Nice work, Dima!")
       return
   }
   sort.Ints(ls)
   // unique ls
   us := ls[:0]
   for i, v := range ls {
       if i == 0 || v != ls[i-1] {
           us = append(us, v)
       }
   }
   answer := 0
   // compute global max rk for pruning
   globalMaxR := 0
   for _, e := range edges {
       if e.r > globalMaxR {
           globalMaxR = e.r
       }
   }
   // for each possible L
   for _, L := range us {
       // prune if even maximum possible loyalty <= current answer
       if globalMaxR - L + 1 <= answer {
           break
       }
       // build graph with edges having l <= L
       adj := make([][]Edge, n+1)
       for _, e := range edges {
           if e.l <= L {
               adj[e.u] = append(adj[e.u], Edge{to: e.v, rk: e.r})
               adj[e.v] = append(adj[e.v], Edge{to: e.u, rk: e.r})
           }
       }
       W := widestPath(n, adj)
       if W >= L {
           loyalty := W - L + 1
           if loyalty > answer {
               answer = loyalty
           }
       }
   }
   if answer <= 0 {
       fmt.Println("Nice work, Dima!")
   } else {
       fmt.Println(answer)
   }
}
