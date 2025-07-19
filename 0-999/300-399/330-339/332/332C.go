package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// E represents a problem with attributes a, b and original index id
type E struct {
   a, b, id int
}

// Node is an element in the min-heap, storing a and its index in sorted slice
type Node struct {
   a, idx int
}

// A MinHeap implements heap.Interface based on Node.a
type MinHeap []Node

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].a < h[j].a }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
   *h = append(*h, x.(Node))
}

func (h *MinHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, p, k int
   fmt.Fscan(in, &n, &p, &k)
   es := make([]E, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &es[i].a, &es[i].b)
       es[i].id = i + 1
   }
   // sort by b asc, then a desc
   sort.Slice(es, func(i, j int) bool {
       if es[i].b != es[j].b {
           return es[i].b < es[j].b
       }
       return es[i].a > es[j].a
   })

   // build min-heap of size k from indices [p-k .. n-1]
   h := &MinHeap{}
   heap.Init(h)
   // initial: indices n-k .. n-1
   for i := n - 1; i >= n-k; i-- {
       heap.Push(h, Node{a: es[i].a, idx: i})
   }
   // consider remaining in [p-k .. n-k-1]
   for i := n - k - 1; i >= p-k; i-- {
       if es[i].a > (*h)[0].a {
           heap.Pop(h)
           heap.Push(h, Node{a: es[i].a, idx: i})
       }
   }
   // extract selected by a
   selectedA := make([]int, 0, k)
   last := n
   for h.Len() > 0 {
       nd := heap.Pop(h).(Node)
       selectedA = append(selectedA, es[nd.idx].id)
       if nd.idx < last {
           last = nd.idx
       }
   }
   // select remaining p-k by highest b from prefix [0..last-1]
   selectedB := make([]int, 0, p-k)
   for i := last - 1; len(selectedB) < p-k; i-- {
       selectedB = append(selectedB, es[i].id)
   }
   // output
   for _, id := range selectedA {
       fmt.Fprint(out, id, " ")
   }
   for _, id := range selectedB {
       fmt.Fprint(out, id, " ")
   }
   fmt.Fprintln(out)
}
