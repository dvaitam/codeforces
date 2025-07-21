package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Event represents either an update or query event in the sweep line
type Event struct {
   x    int   // x-coordinate of event
   typ  int   // 0 = update, 1 = query
   y1   int   // range start in y
   y2   int   // range end in y
   val  int64 // for update: delta v; for query: coefficient (+1 or -1)
   idx  int   // query index
}

// BIT supports range update and range sum query on y-axis
type BIT struct {
   n    int
   bit0 []int64
   bit1 []int64
}

func newBIT(n int) *BIT {
   return &BIT{n: n, bit0: make([]int64, n+5), bit1: make([]int64, n+5)}
}

func (b *BIT) upd(bit []int64, i int, v int64) {
   for ; i <= b.n; i += i & -i {
       bit[i] += v
   }
}

// updateRange adds v to all positions in [l..r]
func (b *BIT) updateRange(l, r int, v int64) {
   if l > r {
       return
   }
   b.upd(b.bit0, l, v)
   b.upd(b.bit0, r+1, -v)
   b.upd(b.bit1, l, v*int64(l-1))
   b.upd(b.bit1, r+1, -v*int64(r))
}

// prefixSum returns sum of [1..i]
func (b *BIT) prefixSum(i int) int64 {
   var s0, s1 int64
   j := i
   for ; j > 0; j -= j & -j {
       s0 += b.bit0[j]
       s1 += b.bit1[j]
   }
   return s0*int64(i) - s1
}

// queryRange returns sum in [l..r]
func (b *BIT) queryRange(l, r int) int64 {
   if l > r {
       return 0
   }
   return b.prefixSum(r) - b.prefixSum(l-1)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m, w int
   fmt.Fscan(in, &n, &m, &w)
   events := make([]Event, 0, w*10)
   qcnt := 0
   for i := 0; i < w; i++ {
       var t int
       fmt.Fscan(in, &t)
       if t == 0 {
           var x1, y1, x2, y2 int
           var v int64
           fmt.Fscan(in, &x1, &y1, &x2, &y2, &v)
           // start update at x1, end at x2+1
           events = append(events, Event{x: x1, typ: 0, y1: y1, y2: y2, val: v})
           events = append(events, Event{x: x2 + 1, typ: 0, y1: y1, y2: y2, val: -v})
       } else {
           var x1, y1, x2, y2 int
           fmt.Fscan(in, &x1, &y1, &x2, &y2)
           id := qcnt
           qcnt++
           // Dima's rectangle [x1..x2] x [y1..y2], coeff +1
           events = append(events, Event{x: x2, typ: 1, y1: y1, y2: y2, val: 1, idx: id})
           events = append(events, Event{x: x1 - 1, typ: 1, y1: y1, y2: y2, val: -1, idx: id})
           // Inna's 4 corner rectangles, coeff -1
           // Top-left [1..x1-1] x [1..y1-1]
           events = append(events, Event{x: x1 - 1, typ: 1, y1: 1, y2: y1 - 1, val: -1, idx: id})
           events = append(events, Event{x: 0, typ: 1, y1: 1, y2: y1 - 1, val: 1, idx: id})
           // Top-right [1..x1-1] x [y2+1..m]
           events = append(events, Event{x: x1 - 1, typ: 1, y1: y2 + 1, y2: m, val: -1, idx: id})
           events = append(events, Event{x: 0, typ: 1, y1: y2 + 1, y2: m, val: 1, idx: id})
           // Bottom-left [x2+1..n] x [1..y1-1]
           events = append(events, Event{x: n, typ: 1, y1: 1, y2: y1 - 1, val: -1, idx: id})
           events = append(events, Event{x: x2, typ: 1, y1: 1, y2: y1 - 1, val: 1, idx: id})
           // Bottom-right [x2+1..n] x [y2+1..m]
           events = append(events, Event{x: n, typ: 1, y1: y2 + 1, y2: m, val: -1, idx: id})
           events = append(events, Event{x: x2, typ: 1, y1: y2 + 1, y2: m, val: 1, idx: id})
       }
   }
   // sort by x, updates before queries at same x
   sort.Slice(events, func(i, j int) bool {
       if events[i].x != events[j].x {
           return events[i].x < events[j].x
       }
       return events[i].typ < events[j].typ
   })
   bit := newBIT(m + 2)
   ans := make([]int64, qcnt)
   // sweep line
   for _, e := range events {
       if e.typ == 0 {
           bit.updateRange(e.y1, e.y2, e.val)
       } else {
           res := bit.queryRange(e.y1, e.y2)
           ans[e.idx] += e.val * res
       }
   }
   // output answers
   for i := 0; i < qcnt; i++ {
       fmt.Fprintln(out, ans[i])
   }
}
