package main

import (
   "bufio"
   "fmt"
   "os"
)

// A simple min-heap for int64
type Heap struct {
   data []int64
}
func (h *Heap) Len() int { return len(h.data) }
func (h *Heap) Push(x int64) {
   h.data = append(h.data, x)
   i := len(h.data) - 1
   // sift up
   for i > 0 {
       p := (i - 1) / 2
       if h.data[p] <= h.data[i] {
           break
       }
       h.data[p], h.data[i] = h.data[i], h.data[p]
       i = p
   }
}
func (h *Heap) Pop() int64 {
   n := len(h.data)
   if n == 0 {
       return 0
   }
   ret := h.data[0]
   last := h.data[n-1]
   h.data = h.data[:n-1]
   if n-1 > 0 {
       h.data[0] = last
       // sift down
       i := 0
       for {
           left := 2*i + 1
           if left >= len(h.data) {
               break
           }
           minChild := left
           right := left + 1
           if right < len(h.data) && h.data[right] < h.data[left] {
               minChild = right
           }
           if h.data[minChild] >= h.data[i] {
               break
           }
           h.data[i], h.data[minChild] = h.data[minChild], h.data[i]
           i = minChild
       }
   }
   return ret
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   cnt := make(map[int64]int64, n)
   var ai int64
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &ai)
       cnt[ai]++
   }
   // heap of positions needing carry processing
   h := &Heap{}
   inHeap := make(map[int64]bool)
   for k, v := range cnt {
       if v >= 2 {
           h.Push(k)
           inHeap[k] = true
       }
   }
   // process carries
   for h.Len() > 0 {
       k := h.Pop()
       inHeap[k] = false
       v := cnt[k]
       if v < 2 {
           continue
       }
       carry := v >> 1
       cnt[k] = v & 1
       nk := k + 1
       cnt[nk] += carry
       if cnt[nk] >= 2 && !inHeap[nk] {
           h.Push(nk)
           inHeap[nk] = true
       }
   }
   // compute highest bit and popcount
   var H int64 = -1
   var ones int64
   for k, v := range cnt {
       if v > 0 {
           if k > H {
               H = k
           }
           ones++
       }
   }
   // answer = (H+1) - ones
   // H>=0 since n>=1
   ans := (H + 1) - ones
   fmt.Println(ans)
}
