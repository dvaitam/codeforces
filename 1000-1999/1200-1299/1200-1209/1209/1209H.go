package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "math"
   "os"
)

// Item represents an element in the max-heap
type Item struct {
   v   float64
   idx int
}

// MaxHeap implements a max-heap of Items
type MaxHeap []Item

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].v > h[j].v }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   it := old[n-1]
   *h = old[:n-1]
   return it
}

// Peek returns the top element without removing it
func (h *MaxHeap) Peek() Item {
   return (*h)[0]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var L int
   if _, err := fmt.Fscan(reader, &n, &L); err != nil {
       return
   }
   x0 := make([]int, n)
   y0 := make([]int, n)
   v0 := make([]float64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x0[i], &y0[i], &v0[i])
   }
   var l, v []float64
   if x0[0] > 0 {
       l = append(l, float64(x0[0]))
       v = append(v, 0)
   }
   for i := 0; i < n; i++ {
       l = append(l, float64(y0[i]-x0[i]))
       v = append(v, v0[i])
       if i+1 < n && y0[i] < x0[i+1] {
           l = append(l, float64(x0[i+1]-y0[i]))
           v = append(v, 0)
       }
       if i+1 == n && y0[i] < L {
           l = append(l, float64(L-y0[i]))
           v = append(v, 0)
       }
   }
   m := len(v)
   x := make([]float64, m)
   for i := range x {
       x[i] = 2.0
   }
   t := make([]float64, m)
   h := &MaxHeap{}
   heap.Init(h)
   var bal, res float64
   const eps = 1e-10
   for i := 0; i < m; i++ {
       heap.Push(h, Item{v: v[i], idx: i})
       t[i] = l[i] / (v[i] + x[i])
       bal += t[i] * (1 - x[i])
       res += t[i]
       for bal < -eps {
           if h.Len() == 0 {
               break
           }
           top := h.Peek()
           j := top.idx
           bal -= t[j] * (1 - x[j])
           res -= t[j]
           var newx float64
           if v[j] > eps && bal+l[j]/v[j] < eps {
               newx = 0
           } else {
               newx = (l[j] + bal*v[j]) / (l[j] - bal)
               if newx < 0 {
                   newx = 0
               }
               if newx > 2 {
                   newx = 2
               }
           }
           x[j] = newx
           if x[j] < eps {
               heap.Pop(h)
           }
           t[j] = l[j] / (v[j] + x[j])
           bal += t[j] * (1 - x[j])
           res += t[j]
       }
   }
   fmt.Printf("%.15f\n", res)
}
