package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

// Segment Tree Beats supporting range add and range chmax (max-floor) with range sum
type SegTree struct {
   n        int
   sum      []int64
   minv     []int   // minimum value in node
   sminv    []int   // second minimum
   cntmin   []int   // count of minimum
   addv     []int   // lazy add
}

func NewSegTree(n int) *SegTree {
   size := 1
   for size < n*4 {
       size <<= 1
   }
   st := &SegTree{
       n:      n,
       sum:    make([]int64, size),
       minv:   make([]int, size),
       sminv:  make([]int, size),
       cntmin: make([]int, size),
       addv:   make([]int, size),
   }
   st.build(1, 1, n)
   return st
}

func (st *SegTree) build(o, l, r int) {
   st.addv[o] = 0
   st.sum[o] = 0
   st.minv[o] = 0
   st.sminv[o] = int(1e9)
   st.cntmin[o] = r - l + 1
   if l == r {
       return
   }
   m := (l + r) >> 1
   st.build(o<<1, l, m)
   st.build(o<<1|1, m+1, r)
}

func (st *SegTree) pushAdd(o, l, r, v int) {
   st.sum[o] += int64(v) * int64(r-l+1)
   st.minv[o] += v
   // sminv: if > inf/2, keep as is
   if st.sminv[o] < int(1e9) {
       st.sminv[o] += v
   }
   st.addv[o] += v
}

func (st *SegTree) pushChmax(o int, x int) {
   // assume minv[o] < x < sminv[o]
   st.sum[o] += int64(x - st.minv[o]) * int64(st.cntmin[o])
   st.minv[o] = x
}

func (st *SegTree) pushDown(o, l, r int) {
   if l == r {
       st.addv[o] = 0
       return
   }
   lc, rc := o<<1, o<<1|1
   // propagate add
   if st.addv[o] != 0 {
       m := (l + r) >> 1
       st.pushAdd(lc, l, m, st.addv[o])
       st.pushAdd(rc, m+1, r, st.addv[o])
       st.addv[o] = 0
   }
   // propagate chmax on min
   if st.minv[o] > st.minv[lc] {
       st.pushChmax(lc, st.minv[o])
   }
   if st.minv[o] > st.minv[rc] {
       st.pushChmax(rc, st.minv[o])
   }
}

func (st *SegTree) pushUp(o int) {
   lc, rc := o<<1, o<<1|1
   st.sum[o] = st.sum[lc] + st.sum[rc]
   // merge min
   if st.minv[lc] < st.minv[rc] {
       st.minv[o] = st.minv[lc]
       st.cntmin[o] = st.cntmin[lc]
       st.sminv[o] = min(st.sminv[lc], st.minv[rc])
   } else if st.minv[lc] > st.minv[rc] {
       st.minv[o] = st.minv[rc]
       st.cntmin[o] = st.cntmin[rc]
       st.sminv[o] = min(st.minv[lc], st.sminv[rc])
   } else {
       st.minv[o] = st.minv[lc]
       st.cntmin[o] = st.cntmin[lc] + st.cntmin[rc]
       st.sminv[o] = min(st.sminv[lc], st.sminv[rc])
   }
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

// Range add v on [L,R]
func (st *SegTree) rangeAdd(o, l, r, L, R, v int) {
   if L > R || L > r || R < l {
       return
   }
   if L <= l && r <= R {
       st.pushAdd(o, l, r, v)
       return
   }
   st.pushDown(o, l, r)
   m := (l + r) >> 1
   if L <= m {
       st.rangeAdd(o<<1, l, m, L, R, v)
   }
   if R > m {
       st.rangeAdd(o<<1|1, m+1, r, L, R, v)
   }
   st.pushUp(o)
}

// Range chmax: for all i in [L,R], a[i] = max(a[i], x)
func (st *SegTree) rangeChmax(o, l, r, L, R, x int) {
   if L > R || L > r || R < l || st.minv[o] >= x {
       return
   }
   if L <= l && r <= R && st.sminv[o] > x {
       st.pushChmax(o, x)
       return
   }
   st.pushDown(o, l, r)
   m := (l + r) >> 1
   if L <= m {
       st.rangeChmax(o<<1, l, m, L, R, x)
   }
   if R > m {
       st.rangeChmax(o<<1|1, m+1, r, L, R, x)
   }
   st.pushUp(o)
}

// Range sum query
func (st *SegTree) rangeSum(o, l, r, L, R int) int64 {
   if L > R || L > r || R < l {
       return 0
   }
   if L <= l && r <= R {
       return st.sum[o]
   }
   st.pushDown(o, l, r)
   m := (l + r) >> 1
   var res int64
   if L <= m {
       res += st.rangeSum(o<<1, l, m, L, R)
   }
   if R > m {
       res += st.rangeSum(o<<1|1, m+1, r, L, R)
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, _ := reader.ReadString('\n')
   n, _ := strconv.Atoi(line[:len(line)-1])
   sline, _ := reader.ReadString('\n')
   s := sline[:n]
   st := NewSegTree(n)
   var total int64
   end := 0
   for i := 1; i <= n; i++ {
       if s[i-1] == '1' {
           end++
           k := end
           l := i - k
           // set arr[i] = 1 via chmax to 1 on [1..i]
           st.rangeChmax(1, 1, n, 1, i, 1)
           if k > 1 {
               // add +1 on suffix [l+1..i-1]
               st.rangeAdd(1, 1, n, l+1, i-1, 1)
           }
           if l >= 1 {
               st.rangeChmax(1, 1, n, 1, l, k)
           }
      } else {
          end = 0
      }
       total += st.rangeSum(1, 1, n, 1, i)
   }
   fmt.Println(total)
}
