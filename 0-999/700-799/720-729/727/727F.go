package main

import (
   "bufio"
   "fmt"
   "os"
)

// simple min-heap for int64
type minHeap struct {
   data []int64
}
func (h *minHeap) Init(c int) {
   h.data = h.data[:0]
   if c > cap(h.data) {
       h.data = make([]int64, 0, c)
   }
}
func (h *minHeap) Push(v int64) {
   h.data = append(h.data, v)
   i := len(h.data) - 1
   for i > 0 {
       p := (i - 1) >> 1
       if h.data[p] <= h.data[i] {
           break
       }
       h.data[p], h.data[i] = h.data[i], h.data[p]
       i = p
   }
}
func (h *minHeap) Pop() int64 {
   n := len(h.data)
   v := h.data[0]
   last := h.data[n-1]
   h.data = h.data[:n-1]
   if n-1 > 0 {
       h.data[0] = last
       // sift down
       i := 0
       for {
           l := 2*i + 1
           if l >= len(h.data) {
               break
           }
           r := l + 1
           j := l
           if r < len(h.data) && h.data[r] < h.data[l] {
               j = r
           }
           if h.data[j] >= h.data[i] {
               break
           }
           h.data[i], h.data[j] = h.data[j], h.data[i]
           i = j
       }
   }
   return v
}

// count removals needed with initial mood q, stop early if >limit
func countRemovals(a []int64, q int64, limit int, heap *minHeap) int {
   sum := q
   rem := 0
   heap.Init(len(a))
   for _, v := range a {
       sum += v
       if v < 0 {
           heap.Push(v)
           if sum < 0 {
               // remove most negative
               x := heap.Pop()
               sum -= x
               rem++
               if rem > limit {
                   return rem
               }
           }
       }
   }
   return rem
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m int
   fmt.Fscan(in, &n, &m)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   b := make([]int64, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &b[i])
   }
   // f[k] = minimal q s.t. countRemovals(q) <= k
   // first compute f[0] = max prefix deficit
   var pref int64
   minPref := int64(0)
   for _, v := range a {
       pref += v
       if pref < minPref {
           minPref = pref
       }
   }
   f := make([]int64, n+1)
   if minPref < 0 {
       f[0] = -minPref
   } else {
       f[0] = 0
   }
   heap := new(minHeap)
   // compute for k >=1
   for k := 1; k <= n; k++ {
       if f[k-1] == 0 {
           f[k] = 0
           continue
       }
       lo, hi := int64(0), f[k-1]
       for lo < hi {
           mid := (lo + hi) >> 1
           if countRemovals(a, mid, k, heap) <= k {
               hi = mid
           } else {
               lo = mid + 1
           }
       }
       f[k] = lo
   }
   // answer queries
   for _, q := range b {
       // find minimal k s.t. f[k] <= q
       lo, hi := 0, n
       for lo < hi {
           mid := (lo + hi) >> 1
           if f[mid] <= q {
               hi = mid
           } else {
               lo = mid + 1
           }
       }
       if lo > n {
           lo = n
       }
       fmt.Fprintln(out, lo)
   }
}
