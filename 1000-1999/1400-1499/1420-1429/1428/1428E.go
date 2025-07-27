package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// item holds a carrot of length a split into t parts, with potential gain from one more split
type item struct {
   a     int64
   t     int64
   delta int64
}

// maxHeap implements heap.Interface as a max-heap based on delta
type maxHeap []*item

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].delta > h[j].delta }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) {
   *h = append(*h, x.(*item))
}
func (h *maxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

// cost returns the total cost of splitting a into t parts
func cost(a, t int64) int64 {
   q := a / t
   r := a % t
   // r parts of size q+1, t-r parts of size q
   return r*(q+1)*(q+1) + (t-r)*q*q
}

// gain returns reduction in cost by increasing splits from t to t+1
func gain(a, t int64) int64 {
   return cost(a, t) - cost(a, t+1)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // initial total cost: each carrot as single piece
   var total int64
   for i := 0; i < n; i++ {
       total += a[i] * a[i]
   }
   // build max-heap of initial gains
   h := &maxHeap{}
   heap.Init(h)
   for i := 0; i < n; i++ {
       if a[i] > 1 {
           itm := &item{a: a[i], t: 1, delta: gain(a[i], 1)}
           heap.Push(h, itm)
       }
   }
   // perform k-n additional splits
   extra := k - n
   for i := 0; i < extra; i++ {
       if h.Len() == 0 {
           break
       }
       top := heap.Pop(h).(*item)
       if top.delta <= 0 {
           break
       }
       total -= top.delta
       top.t++
       // push back if can split further
       if top.t < top.a {
           top.delta = gain(top.a, top.t)
           heap.Push(h, top)
       }
   }
   fmt.Fprintln(writer, total)
}
