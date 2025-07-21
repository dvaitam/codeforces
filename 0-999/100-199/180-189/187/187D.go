package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Interval represents a piecewise function segment for offset transformation
type Interval struct {
   start    int64 // inclusive start of x
   a, b     int64 // time_increase = a*x + b
   c, d     int64 // new_offset    = c*x + d
}

// Node holds intervals for a segment tree node
type Node struct {
   ivals []Interval
}

// mergeNode composes two nodes: first apply left, then right
func mergeNode(L, R Node) Node {
   var res Node
   var ivs []Interval
   // prepare R starts
   Rs := R.ivals
   Rstarts := make([]int64, len(Rs))
   for i := range Rs {
       Rstarts[i] = Rs[i].start
   }
   // process each interval in L
   for i, Li := range L.ivals {
       l0 := Li.start
       var l1 int64
       if i+1 < len(L.ivals) {
           l1 = L.ivals[i+1].start
       } else {
           l1 = cycle
       }
       if Li.c == 0 {
           // x1 = d constant
           x1 := Li.d
           // find R interval
           j := sort.Search(len(Rstarts), func(k int) bool { return Rstarts[k] > x1 }) - 1
           if j < 0 {
               j = 0
           }
           Rj := Rs[j]
           // compose
           a := Li.a
           b := Li.b + Rj.b + Rj.a*Li.d
           c := Rj.c * 0
           d := Rj.c*Li.d + Rj.d
           ivs = append(ivs, Interval{start: l0, a: a, b: b, c: c, d: d})
       } else {
           // Li.c == 1: x1 = (x + d) mod cycle
           // prepare R interval ends
           Rends := make([]int64, len(Rstarts))
           for j := range Rstarts {
               if j+1 < len(Rstarts) {
                   Rends[j] = Rstarts[j+1]
               } else {
                   Rends[j] = cycle
               }
           }
           // collect mapped intervals
           type segInfo struct{u, v int64; j int}
           segs := make([]segInfo, 0)
           for j, sr := range Rstarts {
               er := Rends[j]
               // two possible shifts: k=0 and k=1
               for k := int64(0); k < 2; k++ {
                   start := sr - Li.d + k*cycle
                   end := er - Li.d + k*cycle
                   if start < l1 && end > l0 {
                       u := start
                       if u < l0 {
                           u = l0
                       }
                       v := end
                       if v > l1 {
                           v = l1
                       }
                       if u < v {
                           segs = append(segs, segInfo{u: u, v: v, j: j})
                       }
                   }
               }
           }
           // sort segments by u
           sort.Slice(segs, func(i, j int) bool { return segs[i].u < segs[j].u })
           // build ivs from segs
           for _, s := range segs {
               a2 := Rs[s.j].a
               b2 := Rs[s.j].b
               c2 := Rs[s.j].c
               d2 := Rs[s.j].d
               a := Li.a + a2
               b := Li.b + a2*Li.d + b2
               c := c2
               d := c2*Li.d + d2
               ivs = append(ivs, Interval{start: s.u, a: a, b: b, c: c, d: d})
           }
       }
   }
   // merge contiguous with same params
   if len(ivs) == 0 {
       res.ivals = ivs
       return res
   }
   merged := make([]Interval, 0, len(ivs))
   cur := ivs[0]
   for i := 1; i < len(ivs); i++ {
       v := ivs[i]
       if v.a == cur.a && v.b == cur.b && v.c == cur.c && v.d == cur.d {
           continue
       }
       merged = append(merged, cur)
       cur = v
   }
   merged = append(merged, cur)
   res.ivals = merged
   return res
}

var (
   cycle int64
   lastSeg int64
   tree []Node
   n int
   gVar int64
)

// build constructs the segment tree over segments 1..n
// buildTree builds segment tree at node idx for segs[l-1..r-1]
func buildTree(idx, l, r int, segs []int64, g int64) {
   if l == r {
       li := segs[l-1]
       // threshold: x + li < g => x < g - li
       thr := g - li
       var ivs []Interval
       if thr <= 0 {
           // always wait
           ivs = []Interval{{start: 0, a: -1, b: cycle, c: 0, d: 0}}
       } else if thr >= cycle {
           // never wait
           ivs = []Interval{{start: 0, a: 0, b: li, c: 1, d: li}}
       } else {
           ivs = []Interval{
               {start: 0, a: 0, b: li, c: 1, d: li},
               {start: thr, a: -1, b: cycle, c: 0, d: 0},
           }
       }
       tree[idx].ivals = ivs
       return
   }
   mid := (l + r) >> 1
   buildTree(idx<<1, l, mid, segs, g)
   buildTree(idx<<1|1, mid+1, r, segs, g)
   tree[idx] = mergeNode(tree[idx<<1], tree[idx<<1|1])
}
// main
func main() {
   in := bufio.NewReader(os.Stdin)
   var r int64
   fmt.Fscan(in, &n, &gVar, &r)
   cycle = gVar + r
   segs := make([]int64, n+1)
   for i := 0; i <= n; i++ {
       fmt.Fscan(in, &segs[i])
   }
   lastSeg = segs[n]
   tree = make([]Node, 4*n)
   if n > 0 {
       buildTree(1, 1, n, segs[:n], gVar)
   }
   var q int
   fmt.Fscan(in, &q)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i := 0; i < q; i++ {
       var t0 int64
       fmt.Fscan(in, &t0)
       x0 := t0 % cycle
       // apply tree[1]
       iv := tree[1].ivals
       // binary search for interval
       j := sort.Search(len(iv), func(i int) bool { return iv[i].start > x0 }) - 1
       if j < 0 {
           j = 0
       }
       seg := iv[j]
       timeInc := seg.a*x0 + seg.b
       res := t0 + timeInc + lastSeg
       fmt.Fprintln(out, res)
   }
}
