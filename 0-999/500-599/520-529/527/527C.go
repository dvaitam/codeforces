package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// MaxHeap is a max-heap of ints.
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

// BIT is a Fenwick tree for sum queries.
type BIT struct {
   n   int
   bit []int
}

func NewBIT(n int) *BIT {
   return &BIT{n: n, bit: make([]int, n+1)}
}

// Add adds v at index i (1-based).
func (f *BIT) Add(i, v int) {
   for ; i <= f.n; i += i & -i {
       f.bit[i] += v
   }
}

// Sum returns sum of [1..i].
func (f *BIT) Sum(i int) int {
   s := 0
   for ; i > 0; i -= i & -i {
       s += f.bit[i]
   }
   return s
}

// FindByOrder returns the smallest i such that Sum(i) >= k.
// assumes k >= 1 and k <= total sum.
func (f *BIT) FindByOrder(k int) int {
   idx := 0
   bitMask := 1 << uint(0)
   for bitMask<<1 <= f.n {
       bitMask <<= 1
   }
   for mask := bitMask; mask > 0; mask >>= 1 {
       nxt := idx + mask
       if nxt <= f.n && f.bit[nxt] < k {
           idx = nxt
           k -= f.bit[nxt]
       }
   }
   return idx + 1
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var w, h, n int
   fmt.Fscan(in, &w, &h, &n)
   events := make([]struct{ t byte; x int }, n)
   xs := []int{0, w}
   ys := []int{0, h}
   for i := 0; i < n; i++ {
       var typ byte
       var v int
       fmt.Fscan(in, &typ, &v)
       events[i] = struct{ t byte; x int }{typ, v}
       if typ == 'V' {
           xs = append(xs, v)
       } else {
           ys = append(ys, v)
       }
   }
   sort.Ints(xs)
   sort.Ints(ys)
   // unique
   ux := xs[:1]
   for _, v := range xs[1:] {
       if v != ux[len(ux)-1] {
           ux = append(ux, v)
       }
   }
   uy := ys[:1]
   for _, v := range ys[1:] {
       if v != uy[len(uy)-1] {
           uy = append(uy, v)
       }
   }
   xs = ux
   ys = uy
   mx := len(xs)
   my := len(ys)
   idxX := make(map[int]int, mx)
   for i, v := range xs {
       idxX[v] = i
   }
   idxY := make(map[int]int, my)
   for i, v := range ys {
       idxY[v] = i
   }
   bitX := NewBIT(mx)
   bitY := NewBIT(my)
   activeX := make([]bool, mx)
   activeY := make([]bool, my)
   // initial edges
   activeX[0], activeX[mx-1] = true, true
   bitX.Add(1, 1)
   bitX.Add(mx, 1)
   activeY[0], activeY[my-1] = true, true
   bitY.Add(1, 1)
   bitY.Add(my, 1)
   // heaps and removal maps
   hx := &MaxHeap{xs[mx-1] - xs[0]}
   hy := &MaxHeap{ys[my-1] - ys[0]}
   heap.Init(hx)
   heap.Init(hy)
   remX := make(map[int]int)
   remY := make(map[int]int)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for _, e := range events {
       if e.t == 'V' {
           i := idxX[e.x]
           pi := bitX.Sum(i + 1)
           // predecessor at order pi
           predPos := bitX.FindByOrder(pi)
           // successor at order pi+1
           succPos := bitX.FindByOrder(pi + 1)
           predIdx := predPos - 1
           succIdx := succPos - 1
           // old length
           old := xs[succIdx] - xs[predIdx]
           remX[old]++
           // new lengths
           l1 := e.x - xs[predIdx]
           l2 := xs[succIdx] - e.x
           heap.Push(hx, l1)
           heap.Push(hx, l2)
           // activate
           activeX[i] = true
           bitX.Add(i+1, 1)
       } else {
           i := idxY[e.x]
           pi := bitY.Sum(i + 1)
           predPos := bitY.FindByOrder(pi)
           succPos := bitY.FindByOrder(pi + 1)
           predIdx := predPos - 1
           succIdx := succPos - 1
           old := ys[succIdx] - ys[predIdx]
           remY[old]++
           l1 := e.x - ys[predIdx]
           l2 := ys[succIdx] - e.x
           heap.Push(hy, l1)
           heap.Push(hy, l2)
           activeY[i] = true
           bitY.Add(i+1, 1)
       }
       // get max
       for hx.Len() > 0 {
           top := (*hx)[0]
           if remX[top] > 0 {
               remX[top]--
               heap.Pop(hx)
           } else {
               break
           }
       }
       for hy.Len() > 0 {
           top := (*hy)[0]
           if remY[top] > 0 {
               remY[top]--
               heap.Pop(hy)
           } else {
               break
           }
       }
       maxW := (*hx)[0]
       maxH := (*hy)[0]
       fmt.Fprintln(out, int64(maxW)*int64(maxH))
   }
}
