package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// maxHeap implements a max-heap of ints
type maxHeap []int

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *maxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

func can(mid int64, a []int, k, x int) bool {
   var del int
   var sum int64
   h := &maxHeap{}
   heap.Init(h)
   for _, v := range a {
       sum += int64(v)
       heap.Push(h, v)
       // remove items until heap size <= x and sum <= mid
       for h.Len() > x || sum > mid {
           // pop largest
           u := heap.Pop(h).(int)
           sum -= int64(u)
           del++
           if del > k {
               return false
           }
       }
       if h.Len() == x {
           // complete block, reset
           *h = (*h)[:0]
           sum = 0
       }
   }
   // tail heap has size <= x and sum <= mid
   return del <= k
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k, x int
   if _, err := fmt.Fscan(in, &n, &k, &x); err != nil {
       return
   }
   a := make([]int, n)
   var total int64
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
       total += int64(a[i])
   }
   lo, hi := int64(0), total
   var ans int64 = hi
   for lo <= hi {
       mid := (lo + hi) / 2
       if can(mid, a, k, x) {
           ans = mid
           hi = mid - 1
       } else {
           lo = mid + 1
       }
   }
   fmt.Println(ans)
}
