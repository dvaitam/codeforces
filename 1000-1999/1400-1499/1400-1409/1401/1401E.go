package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Fenwick for range update, point query
type Fenwick struct {
   n int
   t []int
}

func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, t: make([]int, n+1)}
}

// add v at index i
func (f *Fenwick) add(i, v int) {
   for ; i <= f.n; i += i & -i {
       f.t[i] += v
   }
}

// range update [l..r] add v
func (f *Fenwick) rangeAdd(l, r, v int) {
   if l > r {
       return
   }
   f.add(l, v)
   if r+1 <= f.n {
       f.add(r+1, -v)
   }
}

// point query at i
func (f *Fenwick) query(i int) int {
   s := 0
   for ; i > 0; i -= i & -i {
       s += f.t[i]
   }
   return s
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   type Hor struct{ y, l, r int }
   hors := make([]Hor, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &hors[i].y, &hors[i].l, &hors[i].r)
   }
   type Event struct{ y, x, c int }
   events := make([]Event, 0, m*2)
   for i := 0; i < m; i++ {
       var x, ly, ry int
       fmt.Fscan(in, &x, &ly, &ry)
       // include horizontals with y in [ly, ry]
       events = append(events, Event{y: ry, x: x, c: 1})
       // subtract horizontals with y < ly -> at y = ly-1
       events = append(events, Event{y: ly - 1, x: x, c: -1})
   }
   sort.Slice(hors, func(i, j int) bool { return hors[i].y < hors[j].y })
   sort.Slice(events, func(i, j int) bool { return events[i].y < events[j].y })

   const MAXX = 1000000
   // Fenwick index from 1..MAXX+2
   fw := NewFenwick(MAXX + 2)
   var inter int64
   pi := 0
   // process events in increasing y
   for _, e := range events {
       // add horizontals with y <= e.y
       for pi < len(hors) && hors[pi].y <= e.y {
           // update range [l..r]
           l := hors[pi].l + 1
           r := hors[pi].r + 1
           fw.rangeAdd(l, r, 1)
           pi++
       }
       // query at x
       qx := e.x + 1
       if qx < 1 {
           qx = 1
       }
       if qx > fw.n {
           qx = fw.n
       }
       cnt := fw.query(qx)
       inter += int64(e.c) * int64(cnt)
   }
   // result = 1 + n + m + intersections
   res := int64(1) + int64(n) + int64(m) + inter
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, res)
}
