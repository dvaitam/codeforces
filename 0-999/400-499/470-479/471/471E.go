package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Disjoint set union with sum of lengths
type DSU struct {
   p    []int
   sum  []int64
}

func NewDSU(n int, lengths []int64) *DSU {
   p := make([]int, n)
   sum := make([]int64, n)
   for i := 0; i < n; i++ {
       p[i] = i
       sum[i] = lengths[i]
   }
   return &DSU{p: p, sum: sum}
}

func (d *DSU) Find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.Find(d.p[x])
   }
   return d.p[x]
}

func (d *DSU) Union(a, b int) {
   ra := d.Find(a)
   rb := d.Find(b)
   if ra == rb {
       return
   }
   // union rb into ra
   d.p[rb] = ra
   d.sum[ra] += d.sum[rb]
}

// segment tree for active horizontals
type SegTree struct {
   n   int
   id  []int
}

func NewSegTree(n int) *SegTree {
   size := 1
   for size < n {
       size <<= 1
   }
   id := make([]int, size*2)
   for i := range id {
       id[i] = -1
   }
   return &SegTree{n: size, id: id}
}

// update position i (0-based) to v (segment id or -1)
func (st *SegTree) Update(i int, v int) {
   i += st.n
   st.id[i] = v
   for i >>= 1; i > 0; i >>= 1 {
       if st.id[i<<1] != -1 {
           st.id[i] = st.id[i<<1]
       } else {
           st.id[i] = st.id[i<<1|1]
       }
   }
}

// collect all active ids in range [l, r], appending to res
func (st *SegTree) collectAll(u, l, r, ql, qr int, res *[]int) {
   if ql > r || qr < l || st.id[u] == -1 {
       return
   }
   if l == r {
       *res = append(*res, st.id[u])
       return
   }
   mid := (l + r) >> 1
   st.collectAll(u<<1, l, mid, ql, qr, res)
   st.collectAll(u<<1|1, mid+1, r, ql, qr, res)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   type Hseg struct{ x1, x2, y int; id int }
   type Vseg struct{ x, y1, y2 int; id int }
   hs := make([]Hseg, 0, n)
   vs := make([]Vseg, 0, n)
   lengths := make([]int64, n)
   for i := 0; i < n; i++ {
       var x1, y1, x2, y2 int
       fmt.Fscan(in, &x1, &y1, &x2, &y2)
       lengths[i] = int64(x2-x1) + int64(y2-y1)
       if y1 == y2 {
           // horizontal
           hs = append(hs, Hseg{x1, x2, y1, i})
       } else {
           // vertical
           vs = append(vs, Vseg{x1, y1, y2, i})
       }
   }
   dsu := NewDSU(n, lengths)
   // collect unique y of horizontals
   ys := make([]int, len(hs))
   for i, h := range hs {
       ys[i] = h.y
   }
   sort.Ints(ys)
   ys = uniqueInts(ys)
   // map y to index
   yIndex := make(map[int]int, len(ys))
   for i, y := range ys {
       yIndex[y] = i
   }
   // events: x, type, id, for hor add/remove store y index, for vert store y1,y2 idx
   const (
       evAdd = 0
       evQuery = 1
       evRem = 2
   )
   type Event struct{ x, typ, id, yi, yj int }
   evs := make([]Event, 0, len(hs)*2+len(vs))
   for _, h := range hs {
       yi := yIndex[h.y]
       evs = append(evs, Event{h.x1, evAdd, h.id, yi, 0})
       evs = append(evs, Event{h.x2, evRem, h.id, yi, 0})
   }
   for _, v := range vs {
       // find y-range in ys
       // first index >= y1, last index <= y2
       yl := sort.SearchInts(ys, v.y1)
       yr := sort.Search(len(ys), func(i int) bool { return ys[i] > v.y2 }) - 1
       if yl <= yr {
           evs = append(evs, Event{v.x, evQuery, v.id, yl, yr})
       }
   }
   sort.Slice(evs, func(i, j int) bool {
       if evs[i].x != evs[j].x {
           return evs[i].x < evs[j].x
       }
       return evs[i].typ < evs[j].typ
   })
   st := NewSegTree(len(ys))
   // process
   for _, e := range evs {
       switch e.typ {
       case evAdd:
           st.Update(e.yi, e.id)
       case evRem:
           st.Update(e.yi, -1)
       case evQuery:
           // collect all horizontals intersecting this vertical
           var found []int
           st.collectAll(1, 0, st.n-1, e.yi, e.yj, &found)
           for _, hid := range found {
               dsu.Union(e.id, hid)
           }
       }
   }
   // find max sum
   var ans int64
   for i := 0; i < n; i++ {
       if dsu.p[i] == i && dsu.sum[i] > ans {
           ans = dsu.sum[i]
       }
   }
   fmt.Fprintln(out, ans)
}

func uniqueInts(a []int) []int {
   j := 0
   for i := 0; i < len(a); i++ {
       if i == 0 || a[i] != a[i-1] {
           a[j] = a[i]
           j++
       }
   }
   return a[:j]
}
