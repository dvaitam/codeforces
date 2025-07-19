package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
   "strconv"
)

// Fast input
type FastReader struct {
   r *bufio.Reader
}
func NewReader() *FastReader {
   return &FastReader{r: bufio.NewReader(os.Stdin)}
}
func (fr *FastReader) ReadInt() int {
   var x int
   var neg bool
   b, _ := fr.r.ReadByte()
   for (b < '0' || b > '9') && b != '-' {
       b, _ = fr.r.ReadByte()
   }
   if b == '-' {
       neg = true
       b, _ = fr.r.ReadByte()
   }
   for b >= '0' && b <= '9' {
       x = x*10 + int(b-'0')
       b, _ = fr.r.ReadByte()
   }
   if neg {
       return -x
   }
   return x
}
func (fr *FastReader) ReadInt64() int64 {
   var x int64
   var neg bool
   b, _ := fr.r.ReadByte()
   for (b < '0' || b > '9') && b != '-' {
       b, _ = fr.r.ReadByte()
   }
   if b == '-' {
       neg = true
       b, _ = fr.r.ReadByte()
   }
   for b >= '0' && b <= '9' {
       x = x*10 + int64(b-'0')
       b, _ = fr.r.ReadByte()
   }
   if neg {
       return -x
   }
   return x
}

// Request (requirement)
type Req struct {
   val, idx int
}
// Shooter with power and cost
type Shooter struct {
   power, cost, idx int
}

// Min-heap of ints
type IntHeap []int
func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) {
   *h = append(*h, x.(int))
}
func (h *IntHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

// Pair for cost and index
type Pair struct { cost, idx int }
// Min-heap of Pair by cost
type PairHeap []Pair
func (h PairHeap) Len() int           { return len(h) }
func (h PairHeap) Less(i, j int) bool { return h[i].cost < h[j].cost }
func (h PairHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *PairHeap) Push(x interface{}) {
   *h = append(*h, x.(Pair))
}
func (h *PairHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

func main() {
   fr := NewReader()
   n := fr.ReadInt()
   m := fr.ReadInt()
   lim := fr.ReadInt64()
   reqs := make([]Req, m)
   for i := 0; i < m; i++ {
       reqs[i] = Req{val: fr.ReadInt(), idx: i}
   }
   shooters := make([]Shooter, n)
   for i := 0; i < n; i++ {
       shooters[i].power = fr.ReadInt()
       shooters[i].idx = i + 1
   }
   for i := 0; i < n; i++ {
       shooters[i].cost = fr.ReadInt()
   }
   sort.Slice(reqs, func(i, j int) bool { return reqs[i].val < reqs[j].val })
   sort.Slice(shooters, func(i, j int) bool { return shooters[i].power < shooters[j].power })
   bel := make([]int, m)

   // check feasibility for group size k
   chk := func(k int) bool {
       var cst int64
       h := &IntHeap{}
       heap.Init(h)
       sp := n - 1
       for i := m - 1; i >= 0; i -= k {
           // push eligible shooters
           for sp >= 0 && shooters[sp].power >= reqs[i].val {
               heap.Push(h, shooters[sp].cost)
               sp--
           }
           if h.Len() == 0 {
               return false
           }
           c := heap.Pop(h).(int)
           cst += int64(c)
           if cst > lim {
               return false
           }
       }
       return true
   }
   // generate assignment for group size k
   gen := func(k int) bool {
       var cst int64
       h := &PairHeap{}
       heap.Init(h)
       sp := n - 1
       for i := m - 1; i >= 0; i -= k {
           for sp >= 0 && shooters[sp].power >= reqs[i].val {
               heap.Push(h, Pair{cost: shooters[sp].cost, idx: shooters[sp].idx})
               sp--
           }
           if h.Len() == 0 {
               return false
           }
           top := (*h)[0]
           cst += int64(top.cost)
           if cst > lim {
               return false
           }
           // assign this shooter to k requirements
           for j := 0; j < k && i-j >= 0; j++ {
               bel[reqs[i-j].idx] = top.idx
           }
           heap.Pop(h)
       }
       return true
   }

   // binary search k
   l, r := 1, m+1
   for l < r {
       mid := (l + r) >> 1
       if chk(mid) {
           r = mid
       } else {
           l = mid + 1
       }
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   if r > m {
       w.WriteString("NO\n")
       return
   }
   w.WriteString("YES\n")
   gen(l)
   for i := 0; i < m; i++ {
       w.WriteString(strconv.Itoa(bel[i]))
       w.WriteByte(' ')
   }
   w.WriteByte('\n')
}
