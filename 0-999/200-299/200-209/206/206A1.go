package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// SeqHeap implements a min-heap of sequences by current head value
type Seq struct {
   value int64
   id    int
}

type SeqHeap []Seq

func (h SeqHeap) Len() int { return len(h) }
func (h SeqHeap) Less(i, j int) bool { return h[i].value < h[j].value }
func (h SeqHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *SeqHeap) Push(x interface{}) { *h = append(*h, x.(Seq)) }
func (h *SeqHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   k := make([]int, n)
   xi := make([]int64, n)
   yi := make([]int64, n)
   mi := make([]int64, n)
   curr := make([]int64, n)
   pos := make([]int, n)
   total := 0
   for i := 0; i < n; i++ {
       var a1 int64
       fmt.Fscan(in, &k[i], &a1, &xi[i], &yi[i], &mi[i])
       curr[i] = a1
       pos[i] = 1
       total += k[i]
   }
   // initialize heap
   h := &SeqHeap{}
   heap.Init(h)
   for i := 0; i < n; i++ {
       if k[i] > 0 {
           heap.Push(h, Seq{value: curr[i], id: i})
       }
   }
   // prepare result storage if needed
   needPrint := total <= 200000
   var resVals []int64
   var resIDs []int
   if needPrint {
       resVals = make([]int64, 0, total)
       resIDs = make([]int, 0, total)
   }
   var badCount int64
   var prev int64
   first := true
   // merge
   for h.Len() > 0 {
       seq := heap.Pop(h).(Seq)
       v := seq.value
       i := seq.id
       if first {
           prev = v
           first = false
       } else {
           if prev > v {
               badCount++
           }
           prev = v
       }
       if needPrint {
           resVals = append(resVals, v)
           resIDs = append(resIDs, i+1)
       }
       // advance sequence i
       pos[i]++
       if pos[i] < k[i] {
           // generate next value
           next := (curr[i]*xi[i] + yi[i]) % mi[i]
           curr[i] = next
           heap.Push(h, Seq{value: next, id: i})
       }
   }
   // output
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprintln(out, badCount)
   if needPrint {
       for idx := range resVals {
           fmt.Fprintln(out, resVals[idx], resIDs[idx])
       }
   }
}
