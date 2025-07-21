package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// Edge represents a possible marriage edge with weight and bitmasks
type Edge struct {
   mMask uint32
   wMask uint32
   r     int
}

// State represents a node in the search tree
type State struct {
   sum      int
   idx      int
   menMask  uint32
   womMask  uint32
}

// PQ implements a min-heap of *State by sum
type PQ []*State

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool { return pq[i].sum < pq[j].sum }
func (pq PQ) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PQ) Push(x interface{}) {
   *pq = append(*pq, x.(*State))
}

func (pq *PQ) Pop() interface{} {
   old := *pq
   n := len(old)
   x := old[n-1]
   *pq = old[:n-1]
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   var t int64
   if _, err := fmt.Fscan(in, &n, &k, &t); err != nil {
       return
   }
   edges := make([]Edge, k)
   for i := 0; i < k; i++ {
       var h, w, r int
       fmt.Fscan(in, &h, &w, &r)
       h--;
       w--;
       edges[i] = Edge{mMask: 1 << uint(h), wMask: 1 << uint(w), r: r}
   }
   // sort edges by increasing weight
   sort.Slice(edges, func(i, j int) bool { return edges[i].r < edges[j].r })

   // min-heap
   var pq PQ
   heap.Init(&pq)
   // initial state
   heap.Push(&pq, &State{sum: 0, idx: 0, menMask: 0, womMask: 0})
   var cnt int64 = 0
   for pq.Len() > 0 {
       s := heap.Pop(&pq).(*State)
       if s.idx == k {
           cnt++
           if cnt == t {
               fmt.Println(s.sum)
               return
           }
           continue
       }
       // skip edge s.idx
       heap.Push(&pq, &State{sum: s.sum, idx: s.idx + 1, menMask: s.menMask, womMask: s.womMask})
       // include edge if possible
       e := edges[s.idx]
       if (s.menMask&e.mMask) == 0 && (s.womMask&e.wMask) == 0 {
           heap.Push(&pq, &State{sum: s.sum + e.r, idx: s.idx + 1,
               menMask: s.menMask | e.mMask,
               womMask: s.womMask | e.wMask,
           })
       }
   }
}
