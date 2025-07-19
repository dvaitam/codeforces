package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Item represents a segment with slope and unique id
type Item struct {
   slope float64
   id    int
}

// MinHeap implements heap.Interface for Items by increasing slope
type MinHeap []Item
func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].slope < h[j].slope }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MinHeap) Pop() interface{} {
   old := *h
   n := len(old)
   it := old[n-1]
   *h = old[:n-1]
   return it
}

// MaxHeap implements heap.Interface for Items by decreasing slope
type MaxHeap []Item
func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].slope > h[j].slope }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   it := old[n-1]
   *h = old[:n-1]
   return it
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
       T--
       var n int
       var a, b int
       fmt.Fscan(reader, &n, &a, &b)
       p := make([]int, n)
       q := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &p[i])
       }
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &q[i])
       }
       sx := float64(a)
       tx := float64(a)
       sy := float64(b)
       ty := float64(b)
       if a == 0 && b == 0 {
           for i := 0; i < n; i++ {
               fmt.Fprintln(writer, 0)
           }
           continue
       }
       // segment storage
       minH := &MinHeap{}
       maxH := &MaxHeap{}
       heap.Init(minH)
       heap.Init(maxH)
       nextID := 0
       xlenMap := make(map[int]float64)
       const eps = 1e-9
       for i := 0; i < n; i++ {
           pi := float64(p[i])
           qi := float64(q[i])
           sx -= pi
           tx += pi
           sy += qi
           ty -= qi
           // insert new segment
           nextID++
           sid := nextID
           slope := qi / pi
           xlen := pi + pi
           xlenMap[sid] = xlen
           heap.Push(minH, Item{slope: slope, id: sid})
           heap.Push(maxH, Item{slope: slope, id: sid})
           // fix sx < 0 by removing from min slope side
           for sx < -eps {
               // pop smallest slope
               var slope0 float64
               var llen float64
               // get valid segment
               for {
                   it := heap.Pop(minH).(Item)
                   if v, ok := xlenMap[it.id]; ok {
                       slope0 = it.slope
                       llen = v
                       delete(xlenMap, it.id)
                       break
                   }
               }
               if sx+llen > 0 {
                   // partial restore
                   d := sx + llen
                   newLen := d
                   // remaining part
                   rem := llen - d
                   // reinsert remaining
                   nextID++
                   sid2 := nextID
                   xlenMap[sid2] = rem
                   heap.Push(minH, Item{slope: slope0, id: sid2})
                   heap.Push(maxH, Item{slope: slope0, id: sid2})
                   llen = newLen
               }
               sx += llen
               sy -= llen * slope0
           }
           // fix ty < 0 by removing from max slope side
           for ty < -eps {
               var slope0 float64
               var llen float64
               for {
                   it := heap.Pop(maxH).(Item)
                   if v, ok := xlenMap[it.id]; ok {
                       slope0 = it.slope
                       llen = v
                       delete(xlenMap, it.id)
                       break
                   }
               }
               if ty+slope0*llen > 0 {
                   // partial restore
                   d := ty/slope0 + llen
                   newLen := d
                   rem := llen - d
                   nextID++
                   sid2 := nextID
                   xlenMap[sid2] = rem
                   heap.Push(minH, Item{slope: slope0, id: sid2})
                   heap.Push(maxH, Item{slope: slope0, id: sid2})
                   llen = newLen
               }
               tx -= llen
               ty += slope0 * llen
           }
           // output tx
           // precision sufficient for 1e-6
           fmt.Fprintf(writer, "%.10f\n", tx)
       }
   }
}
