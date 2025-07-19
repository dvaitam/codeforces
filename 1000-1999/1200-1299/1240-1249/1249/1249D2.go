package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// Segment represents an interval [x, y] with an id
type Segment struct {
   x, y int
   id   int
}

// maxHeap is a max-heap of segments based on y
type maxHeap []Segment

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].y > h[j].y }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) {
   *h = append(*h, x.(Segment))
}
func (h *maxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(reader, &n, &k)
   segments := make([]Segment, n)
   const maxCoord = 200010
   // coordinate bounds up to maxCoord, allocate extra for difference array
   c := make([]int, maxCoord+2)
   minn, maxx := maxCoord, 0
   for i := 0; i < n; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       segments[i] = Segment{x: x, y: y, id: i + 1}
       if x < minn {
           minn = x
       }
       if y > maxx {
           maxx = y
       }
       if x < len(c) {
           c[x]++
       }
       if y+1 < len(c) {
           c[y+1]--
       }
   }
   // sort segments by starting x
   sort.Slice(segments, func(i, j int) bool {
       return segments[i].x < segments[j].x
   })
   h := &maxHeap{}
   heap.Init(h)
   ans := make([]int, 0, n)
   flag := 0
   // sweep line over coordinates
   for i := minn; i <= maxx; i++ {
       // update prefix sum for counts
       if i > 0 {
           c[i] += c[i-1]
       }
       // push all segments starting at i
       for flag < n && segments[flag].x <= i {
           heap.Push(h, segments[flag])
           flag++
       }
       // while count exceeds k, remove the segment with max y
       for c[i] > k {
           top := heap.Pop(h).(Segment)
           ans = append(ans, top.id)
           // remove its effect from i to its end
           c[i]--
           if top.y+1 < len(c) {
               c[top.y+1]++
           }
       }
   }
   // output result
   sort.Ints(ans)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, len(ans))
   if len(ans) > 0 {
       for i, id := range ans {
           if i > 0 {
               w.WriteString(" ")
           }
           fmt.Fprint(w, id)
       }
       fmt.Fprintln(w)
   }
}
