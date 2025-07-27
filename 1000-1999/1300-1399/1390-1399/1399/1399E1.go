package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

type Edge struct {
   to, id int
   w int64
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   for t > 0 {
       t--
       var n int
       var S int64
       fmt.Fscan(in, &n, &S)
       adj := make([][]Edge, n+1)
       m := n - 1
       weights := make([]int64, m)
       for i := 0; i < m; i++ {
           var u, v int
           var w int64
           fmt.Fscan(in, &u, &v, &w)
           adj[u] = append(adj[u], Edge{to: v, id: i, w: w})
           adj[v] = append(adj[v], Edge{to: u, id: i, w: w})
           weights[i] = w
       }
       leafCount := make([]int64, m)
       // DFS to compute leaf counts
       type Frame struct{v, parent, eid, idx int; leaves int64}
       stack := make([]Frame, 0, n)
       stack = append(stack, Frame{v: 1, parent: 0, eid: -1, idx: 0, leaves: 0})
       for len(stack) > 0 {
           f := &stack[len(stack)-1]
           if f.idx < len(adj[f.v]) {
               e := adj[f.v][f.idx]
               f.idx++
               if e.to == f.parent {
                   continue
               }
               stack = append(stack, Frame{v: e.to, parent: f.v, eid: e.id, idx: 0, leaves: 0})
           } else {
               // process node
               if f.leaves == 0 {
                   f.leaves = 1
               }
               if f.eid >= 0 {
                   leafCount[f.eid] = f.leaves
               }
               count := f.leaves
               stack = stack[:len(stack)-1]
               if len(stack) > 0 {
                   stack[len(stack)-1].leaves += count
               }
           }
       }
       // initial sum and heap
       var total int64
       pq := &MaxHeap{}
       heap.Init(pq)
       for i := 0; i < m; i++ {
           total += weights[i] * leafCount[i]
           // compute initial saving
           saving := (weights[i] - weights[i]/2) * leafCount[i]
           heap.Push(pq, &Item{saving: saving, id: i})
       }
       var moves int64
       for total > S {
           it := heap.Pop(pq).(*Item)
           if it.saving <= 0 {
               break
           }
           total -= it.saving
           moves++
           i := it.id
           weights[i] /= 2
           newSave := (weights[i] - weights[i]/2) * leafCount[i]
           heap.Push(pq, &Item{saving: newSave, id: i})
       }
       fmt.Fprintln(out, moves)
   }
}

// Max-heap of Items
type Item struct {
   saving int64
   id     int
}
// MaxHeap implements heap.Interface
type MaxHeap []*Item
func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].saving > h[j].saving }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) {
   *h = append(*h, x.(*Item))
}
func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   it := old[n-1]
   *h = old[:n-1]
   return it
}
