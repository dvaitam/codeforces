package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 1000000000

// Segment tree beats for range chmin and point query
type SegBeat struct {
   n    int
   mx   []int
   se   []int
   cnt  []int
}

func NewSegBeat(n int) *SegBeat {
   size := 1
   for size < n*4 {
       size <<= 1
   }
   sb := &SegBeat{
       n:   n,
       mx:  make([]int, size),
       se:  make([]int, size),
       cnt: make([]int, size),
   }
   sb.build(1, 1, n)
   return sb
}

func (sb *SegBeat) build(p, l, r int) {
   if l == r {
       sb.mx[p] = INF
       sb.se[p] = -1
       sb.cnt[p] = 1
       return
   }
   m := (l + r) >> 1
   sb.build(p<<1, l, m)
   sb.build(p<<1|1, m+1, r)
   sb.pull(p)
}

func (sb *SegBeat) pull(p int) {
   lch, rch := p<<1, p<<1|1
   if sb.mx[lch] == sb.mx[rch] {
       sb.mx[p] = sb.mx[lch]
       sb.cnt[p] = sb.cnt[lch] + sb.cnt[rch]
       sb.se[p] = max(sb.se[lch], sb.se[rch])
   } else {
       if sb.mx[lch] > sb.mx[rch] {
           sb.mx[p] = sb.mx[lch]
           sb.cnt[p] = sb.cnt[lch]
           sb.se[p] = max(sb.se[lch], sb.mx[rch])
       } else {
           sb.mx[p] = sb.mx[rch]
           sb.cnt[p] = sb.cnt[rch]
           sb.se[p] = max(sb.mx[lch], sb.se[rch])
       }
   }
}

func (sb *SegBeat) applyChmin(p, x int) {
   if sb.mx[p] <= x {
       return
   }
   sb.mx[p] = x
}

func (sb *SegBeat) push(p int) {
   // push cap to children
   sb.applyChmin(p<<1, sb.mx[p])
   sb.applyChmin(p<<1|1, sb.mx[p])
}

// Range chmin: cap values in [L,R] to x
func (sb *SegBeat) RangeChmin(L, R, x, p, l, r int) {
   if r < L || R < l || sb.mx[p] <= x {
       return
   }
   if L <= l && r <= R && sb.se[p] < x {
       sb.applyChmin(p, x)
       return
   }
   sb.push(p)
   m := (l + r) >> 1
   sb.RangeChmin(L, R, x, p<<1, l, m)
   sb.RangeChmin(L, R, x, p<<1|1, m+1, r)
   sb.pull(p)
}

// Point query value at idx
func (sb *SegBeat) PointQuery(idx, p, l, r int) int {
   if l == r {
       return sb.mx[p]
   }
   sb.push(p)
   m := (l + r) >> 1
   if idx <= m {
       return sb.PointQuery(idx, p<<1, l, m)
   }
   return sb.PointQuery(idx, p<<1|1, m+1, r)
}

// Simple segment tree for range max and point assign
type SegMax struct {
   n int
   mx []int
}

func NewSegMax(a []int) *SegMax {
   n := len(a) - 1
   size := 1
   for size < n*4 {
       size <<= 1
   }
   sm := &SegMax{n: n, mx: make([]int, size)}
   sm.build(1, 1, n, a)
   return sm
}

func (sm *SegMax) build(p, l, r int, a []int) {
   if l == r {
       sm.mx[p] = a[l]
       return
   }
   m := (l + r) >> 1
   sm.build(p<<1, l, m, a)
   sm.build(p<<1|1, m+1, r, a)
   sm.mx[p] = max(sm.mx[p<<1], sm.mx[p<<1|1])
}

func (sm *SegMax) PointAssign(idx, v, p, l, r int) {
   if l == r {
       sm.mx[p] = v
       return
   }
   m := (l + r) >> 1
   if idx <= m {
       sm.PointAssign(idx, v, p<<1, l, m)
   } else {
       sm.PointAssign(idx, v, p<<1|1, m+1, r)
   }
   sm.mx[p] = max(sm.mx[p<<1], sm.mx[p<<1|1])
}

func (sm *SegMax) RangeMax(L, R, p, l, r int) int {
   if r < L || R < l {
       return 0
   }
   if L <= l && r <= R {
       return sm.mx[p]
   }
   m := (l + r) >> 1
   return max(sm.RangeMax(L, R, p<<1, l, m), sm.RangeMax(L, R, p<<1|1, m+1, r))
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   type Op struct{ t, l, r, x int }
   ops := make([]Op, m+1)
   firstAssign := make([]int, n+1)
   for i := 1; i <= n; i++ {
       firstAssign[i] = m + 1
   }
   for i := 1; i <= m; i++ {
       fmt.Fscan(reader, &ops[i].t)
       if ops[i].t == 1 {
           fmt.Fscan(reader, &ops[i].l, &ops[i].r, &ops[i].x)
       } else {
           fmt.Fscan(reader, &ops[i].l, &ops[i].x)
           // use l as k, x as d in type 2
           if firstAssign[ops[i].l] > i {
               firstAssign[ops[i].l] = i
           }
       }
   }
   // pos event at time t: query UB for pos i before op t
   posEvent := make([][]int, m+2)
   for i := 1; i <= n; i++ {
       t := firstAssign[i]
       posEvent[t] = append(posEvent[t], i)
   }
   sb := NewSegBeat(n)
   UB := make([]int, n+1)
   // sweep time
   for t := 1; t <= m; t++ {
       // before op t, record UB for any positions with T_i == t
       for _, i := range posEvent[t] {
           UB[i] = sb.PointQuery(i, 1, 1, n)
       }
       if ops[t].t == 1 {
           sb.RangeChmin(ops[t].l, ops[t].r, ops[t].x, 1, 1, n)
       }
   }
   // time m+1
   for _, i := range posEvent[m+1] {
       UB[i] = sb.PointQuery(i, 1, 1, n)
   }
   // initial B
   B := make([]int, n+1)
   for i := 1; i <= n; i++ {
       B[i] = UB[i]
   }
   // simulate
   sm := NewSegMax(B)
   for t := 1; t <= m; t++ {
       if ops[t].t == 1 {
           v := sm.RangeMax(ops[t].l, ops[t].r, 1, 1, n)
           if v != ops[t].x {
               fmt.Fprintln(writer, "NO")
               return
           }
       } else {
           sm.PointAssign(ops[t].l, ops[t].x, 1, 1, n)
       }
   }
   fmt.Fprintln(writer, "YES")
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, B[i])
   }
   fmt.Fprintln(writer)
}
