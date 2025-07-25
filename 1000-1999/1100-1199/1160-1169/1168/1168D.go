package main

import (
   "bufio"
   "fmt"
   "os"
)

// Segment tree for range add and range max
type SegTree struct {
   n    int
   max  []int
   add  []int
}

func NewSegTree(n int) *SegTree {
   size := 1
   for size < n {
       size <<= 1
   }
   max := make([]int, size<<1)
   add := make([]int, size<<1)
   return &SegTree{n: n, max: max, add: add}
}

func (st *SegTree) build(idx, l, r int, arr []int) {
   if l == r {
       st.max[idx] = arr[l]
       return
   }
   mid := (l + r) >> 1
   st.build(idx<<1, l, mid, arr)
   st.build(idx<<1|1, mid+1, r, arr)
   if st.max[idx<<1] > st.max[idx<<1|1] {
       st.max[idx] = st.max[idx<<1]
   } else {
       st.max[idx] = st.max[idx<<1|1]
   }
}

func (st *SegTree) Build(arr []int) {
   st.build(1, 0, st.n-1, arr)
}

func (st *SegTree) apply(idx, v int) {
   st.max[idx] += v
   st.add[idx] += v
}

func (st *SegTree) push(idx int) {
   if st.add[idx] != 0 {
       st.apply(idx<<1, st.add[idx])
       st.apply(idx<<1|1, st.add[idx])
       st.add[idx] = 0
   }
}

func (st *SegTree) update(idx, l, r, ql, qr, v int) {
   if ql <= l && r <= qr {
       st.apply(idx, v)
       return
   }
   st.push(idx)
   mid := (l + r) >> 1
   if ql <= mid {
       st.update(idx<<1, l, mid, ql, qr, v)
   }
   if qr > mid {
       st.update(idx<<1|1, mid+1, r, ql, qr, v)
   }
   // recalc
   if st.max[idx<<1] > st.max[idx<<1|1] {
       st.max[idx] = st.max[idx<<1]
   } else {
       st.max[idx] = st.max[idx<<1|1]
   }
}

// RangeAdd adds v to [l..r]
func (st *SegTree) RangeAdd(l, r, v int) {
   st.update(1, 0, st.n-1, l, r, v)
}

// Max returns max over all
func (st *SegTree) Max() int {
   return st.max[1]
}

var (
   children [][]int
   edgeChar []byte
   lidx, ridx []int
   leafCount int
   leafArr [26][]int
   depthBad bool
   Ddepth int
)

func dfs(v, d int, counts []int) {
   if v != 1 {
       c := edgeChar[v]
       if c != '?' {
           counts[c-'a']++
       }
   }
   if len(children[v]) == 0 {
       // leaf
       if leafCount == 0 {
           Ddepth = d
       } else if d != Ddepth {
           depthBad = true
       }
       lidx[v] = leafCount
       for i := 0; i < 26; i++ {
           leafArr[i][leafCount] = counts[i]
       }
       leafCount++
       ridx[v] = leafCount - 1
   } else {
       lidx[v] = leafCount
       for _, u := range children[v] {
           dfs(u, d+1, counts)
       }
       ridx[v] = leafCount - 1
   }
   if v != 1 {
       c := edgeChar[v]
       if c != '?' {
           counts[c-'a']--
       }
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, q int
   fmt.Fscan(in, &n, &q)
   children = make([][]int, n+1)
   edgeChar = make([]byte, n+1)
   parent := make([]int, n+1)
   for i := 2; i <= n; i++ {
       var p int
       var c byte
       fmt.Fscan(in, &p, &c)
       parent[i] = p
       children[p] = append(children[p], i)
       edgeChar[i] = c
   }
   // prepare leaf arrays
   lidx = make([]int, n+1)
   ridx = make([]int, n+1)
   for i := 0; i < 26; i++ {
       leafArr[i] = make([]int, n) // upper bound
   }
   counts := make([]int, 26)
   dfs(1, 0, counts)
   if depthBad {
       // consume queries and print Fou
       for i := 0; i < q; i++ {
           var v int
           var c byte
           fmt.Fscan(in, &v, &c)
           fmt.Fprintln(out, "Fou")
       }
       return
   }
   L := leafCount
   // build segment trees
   st := make([]*SegTree, 26)
   mval := make([]int, 26)
   var sumM int
   var sumW int64
   for i := 0; i < 26; i++ {
       arr := leafArr[i][:L]
       t := NewSegTree(L)
       t.Build(arr)
       st[i] = t
       m := t.Max()
       mval[i] = m
       sumM += m
       sumW += int64(m * (i + 1))
   }
   const totalInd = 351
   // process queries
   for qi := 0; qi < q; qi++ {
       var v int
       var cb byte
       fmt.Fscan(in, &v, &cb)
       old := edgeChar[v]
       if old != cb {
           // remove old
           if old != '?' {
               ci := int(old - 'a')
               l, r := lidx[v], ridx[v]
               st[ci].RangeAdd(l, r, -1)
               newm := st[ci].Max()
               if newm != mval[ci] {
                   sumM += newm - mval[ci]
                   sumW += int64((newm - mval[ci]) * (ci + 1))
                   mval[ci] = newm
               }
           }
           // add new
           if cb != '?' {
               ci := int(cb - 'a')
               l, r := lidx[v], ridx[v]
               st[ci].RangeAdd(l, r, +1)
               newm := st[ci].Max()
               if newm != mval[ci] {
                   sumM += newm - mval[ci]
                   sumW += int64((newm - mval[ci]) * (ci + 1))
                   mval[ci] = newm
               }
           }
           edgeChar[v] = cb
       }
       if sumM > Ddepth {
           fmt.Fprintln(out, "Fou")
       } else {
           ans := sumW + int64(Ddepth-sumM)*totalInd
           fmt.Fprintf(out, "Shi %d\n", ans)
       }
   }
}
