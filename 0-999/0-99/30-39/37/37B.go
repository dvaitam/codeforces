package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// max-heap for spell damage
type Item struct {
   dmg int
   idx int
}
type MaxHeap []Item
func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].dmg > h[j].dmg }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   it := old[n-1]
   *h = old[:n-1]
   return it
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var N, maxH, reg int
   if _, err := fmt.Fscan(in, &N, &maxH, &reg); err != nil {
       return
   }
   type Spell struct{ thr int; dmg, idx int }
   spells := make([]Spell, N)
   for i := 0; i < N; i++ {
       var powi, dmgi int
       fmt.Fscan(in, &powi, &dmgi)
       // threshold compare: 100*H <= thrPct
       spells[i] = Spell{thr: powi * maxH, dmg: dmgi, idx: i + 1}
   }
   // sort by threshold descending
   sort.Slice(spells, func(i, j int) bool {
       return spells[i].thr > spells[j].thr
   })
   var h int = maxH
   var D int = 0
   actions := make([][2]int, 0, N)
   heapItems := &MaxHeap{}
   heap.Init(heapItems)
   j := 0
   // helper to push available spells
   pushAvail := func() {
       for j < N && spells[j].thr >= 100*h {
           heap.Push(heapItems, Item{dmg: spells[j].dmg, idx: spells[j].idx})
           j++
       }
   }
   // time 0: can use one scroll
   pushAvail()
   if heapItems.Len() == 0 {
       fmt.Println("NO")
       return
   }
   // use best at t=0
   it := heap.Pop(heapItems).(Item)
   D += it.dmg
   actions = append(actions, [2]int{0, it.idx})
   // simulate
   var t int
   for t = 1; ; t++ {
       // damage phase
       h -= D
       // regen phase (always applies)
       h += reg
       if h > maxH {
           h = maxH
       }
       // check defeat at end of second
       if h <= 0 {
           break
       }
       // add newly available spells
       pushAvail()
       if heapItems.Len() > 0 {
           it = heap.Pop(heapItems).(Item)
           D += it.dmg
           actions = append(actions, [2]int{t, it.idx})
       } else if D <= reg {
           fmt.Println("NO")
           return
       }
       // else continue until defeat
   }
   // output result
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, "YES")
   fmt.Fprintf(w, "%d %d\n", t, len(actions))
   for _, a := range actions {
       fmt.Fprintf(w, "%d %d\n", a[0], a[1])
   }
}
