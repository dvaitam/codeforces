package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "math"
   "os"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func min(a, b float64) float64 {
   if a < b {
       return a
   }
   return b
}

// Item holds index and priority
type Item struct {
   idx int
   pri float64
}

// MinHeap implements a min-heap of Items with update support
type MinHeap struct {
   items []Item
   pos   []int
}

func NewMinHeap(n int) *MinHeap {
   h := &MinHeap{
       items: make([]Item, n),
       pos:   make([]int, n),
   }
   for i := 0; i < n; i++ {
       h.items[i] = Item{idx: i, pri: 0}
       h.pos[i] = i
   }
   heap.Init(h)
   return h
}

func (h *MinHeap) Len() int { return len(h.items) }
func (h *MinHeap) Less(i, j int) bool { return h.items[i].pri < h.items[j].pri }
func (h *MinHeap) Swap(i, j int) {
   h.items[i], h.items[j] = h.items[j], h.items[i]
   h.pos[h.items[i].idx] = i
   h.pos[h.items[j].idx] = j
}
func (h *MinHeap) Push(x interface{}) {
   // not used
}
func (h *MinHeap) Pop() interface{} {
   // not used
   return nil
}

// Update sets priority of element idx and fixes heap
func (h *MinHeap) Update(idx int, pri float64) {
   i := h.pos[idx]
   h.items[i].pri = pri
   heap.Fix(h, i)
}

// Top returns the minimal Item
func (h *MinHeap) Top() Item {
   return h.items[0]
}

func readInt() int {
   var x int
   fmt.Fscan(reader, &x)
   return x
}

func main() {
   defer writer.Flush()
   N := readInt()
   T := readInt()
   Q := readInt()
   P := make([]float64, N)
   L := make([]int, N)
   for i := 0; i < N; i++ {
       P[i] = float64(readInt())
   }
   for i := 0; i < N; i++ {
       L[i] = readInt()
   }
   bet := make([]int, N)
   const INF = 1e300
   // initialize heaps
   h1 := NewMinHeap(N)
   h2 := NewMinHeap(N)
   // initial priorities
   for i := 0; i < N; i++ {
       // add priority: -P * (min(0.5,(b+1)/(L+b+1)) - min(0.5,b/(L+b)))
       addVal := P[i] * (min(0.5, 1.0/float64(L[i]+1)))
       pri1 := -addVal
       h1.Update(i, pri1)
       // remove priority
       pri2 := INF
       h2.Update(i, pri2)
   }
   res := 0.0
   // initial allocate T bets
   for t := 0; t < T; t++ {
       it := h1.Top()
       i := it.idx
       benefit := -it.pri
       bet[i]++
       res += benefit
       // update this index
       // new add priority
       b := bet[i]
       li := L[i]
       addVal := P[i] * (min(0.5, float64(b+1)/float64(li+b+1)) - min(0.5, float64(b)/float64(li+b)))
       h1.Update(i, -addVal)
       // new remove priority
       var remPri float64
       if b > 0 {
           remVal := min(0.5, float64(b-1)/float64(li+b-1)) - min(0.5, float64(b)/float64(li+b))
           remPri = -P[i] * remVal
       } else {
           remPri = INF
       }
       h2.Update(i, remPri)
   }
   // process queries
   const eps = 1e-10
   for qi := 0; qi < Q; qi++ {
       x := readInt()
       y := readInt() - 1
       // adjust baseline
       if bet[y] > 0 {
           res -= P[y] * min(0.5, float64(bet[y])/float64(L[y]+bet[y]))
       }
       if x == 1 {
           L[y]++
       } else {
           L[y]--
       }
       if bet[y] > 0 {
           res += P[y] * min(0.5, float64(bet[y])/float64(L[y]+bet[y]))
       }
       // update y
       b := bet[y]
       li := L[y]
       addVal := P[y] * (min(0.5, float64(b+1)/float64(li+b+1)) - min(0.5, float64(b)/float64(li+b)))
       h1.Update(y, -addVal)
       var remPri float64
       if b > 0 {
           remVal := min(0.5, float64(b-1)/float64(li+b-1)) - min(0.5, float64(b)/float64(li+b))
           remPri = -P[y] * remVal
       } else {
           remPri = INF
       }
       h2.Update(y, remPri)
       // rebalance
       for {
           it1 := h1.Top()
           it2 := h2.Top()
           i := it1.idx
           j := it2.idx
           gainAdd := -it1.pri
           gainRem := it2.pri
           if i != j && gainAdd > gainRem+eps {
               // transfer
               bet[i]++
               bet[j]--
               res += gainAdd
               res -= gainRem
               // update i
               bi := bet[i]
               lii := L[i]
               addV := P[i] * (min(0.5, float64(bi+1)/float64(lii+bi+1)) - min(0.5, float64(bi)/float64(lii+bi)))
               h1.Update(i, -addV)
               var remP float64
               if bi > 0 {
                   rv := min(0.5, float64(bi-1)/float64(lii+bi-1)) - min(0.5, float64(bi)/float64(lii+bi))
                   remP = -P[i] * rv
               } else {
                   remP = INF
               }
               h2.Update(i, remP)
               // update j
               bj := bet[j]
               lji := L[j]
               addVj := P[j] * (min(0.5, float64(bj+1)/float64(lji+bj+1)) - min(0.5, float64(bj)/float64(lji+bj)))
               h1.Update(j, -addVj)
               var remPj float64
               if bj > 0 {
                   rvj := min(0.5, float64(bj-1)/float64(lji+bj-1)) - min(0.5, float64(bj)/float64(lji+bj))
                   remPj = -P[j] * rvj
               } else {
                   remPj = INF
               }
               h2.Update(j, remPj)
           } else {
               break
           }
       }
       // output
       fmt.Fprintf(writer, "%.15f\n", res)
   }
}
