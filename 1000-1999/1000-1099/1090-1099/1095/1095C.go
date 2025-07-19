package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// MaxHeap implements a max-heap of ints.
type MaxHeap []int

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   // Count set bits in n
   ones := 0
   tmp := n
   for tmp > 0 {
       ones += tmp & 1
       tmp >>= 1
   }
   if k < ones || k > n {
       fmt.Fprintln(out, "NO")
       return
   }
   fmt.Fprintln(out, "YES")
   // Initialize max-heap with powers of two from binary decomposition
   h := &MaxHeap{}
   heap.Init(h)
   tmp = n
   bit := 0
   for tmp > 0 {
       if tmp&1 == 1 {
           heap.Push(h, 1<<bit)
       }
       tmp >>= 1
       bit++
   }
   // Split until we have exactly k parts
   for h.Len() < k {
       x := heap.Pop(h).(int)
       half := x / 2
       heap.Push(h, half)
       heap.Push(h, half)
   }
   // Output the parts
   for h.Len() > 0 {
       x := heap.Pop(h).(int)
       fmt.Fprint(out, x, " ")
   }
   fmt.Fprintln(out)
}
