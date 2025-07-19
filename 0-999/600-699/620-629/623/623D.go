package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Item holds the state for each probability sequence
type Item struct {
   key  float64 // priority key
   pcur float64 // current cumulative probability
   idx  int     // index in original array
}

// A MaxHeap implements heap.Interface for Items based on key (max-heap)
type MaxHeap []Item

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].key > h[j].key }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   it := old[n-1]
   *h = old[0 : n-1]
   return it
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var N int
   if _, err := fmt.Fscan(reader, &N); err != nil {
       return
   }
   P := make([]int, N)
   p := make([]float64, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(reader, &P[i])
       p[i] = float64(P[i]) / 100.0
   }
   // initialize heap and base probability
   h := &MaxHeap{}
   heap.Init(h)
   prob := 1.0
   for i := 0; i < N; i++ {
       pi := p[i]
       prob *= pi
       // initial pcur is one trial, next pcur after two trials is p2 = 1-(1-pi)^2
       // key = p2 / pcur
       pcur := pi
       next := 1 - (1-pcur)*(1-pi)
       key := next / pcur
       heap.Push(h, Item{key: key, pcur: pcur, idx: i})
   }
   // expected value accumulator
   ans := float64(N) * prob
   // simulate steps
   for step := N + 1; step <= 200000; step++ {
       // pop best item
       it := heap.Pop(h).(Item)
       pcur := it.pcur
       i := it.idx
       prevProb := prob
       // remove old contribution, add new
       prob /= pcur
       // update pcur: now with one more trial
       pcur = 1 - (1-pcur)*(1-p[i])
       prob *= pcur
       // compute next key for further trial
       next := 1 - (1-pcur)*(1-p[i])
       key := next / pcur
       heap.Push(h, Item{key: key, pcur: pcur, idx: i})
       ans += (prob - prevProb) * float64(step)
   }
   fmt.Fprintf(writer, "%.15f\n", ans)
}
