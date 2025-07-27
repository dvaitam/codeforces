package main

import (
   "bufio"
   "fmt"
   "os"
)

type SegTree struct {
   n    int
   tree []int64
}

func NewSegTree(n int) *SegTree {
   size := 1
   for size < n {
       size <<= 1
   }
   return &SegTree{n: size, tree: make([]int64, 2*size)}
}

// Build initializes the tree leaves with vals (1-based idx up to nOrig)
func (st *SegTree) Build(vals []int64) {
   // vals is 1-based, length <= original n
   for i := 1; i < len(vals); i++ {
       st.tree[st.n+i-1] = vals[i]
   }
   for i := st.n - 1; i > 0; i-- {
       st.tree[i] = st.tree[2*i] + st.tree[2*i+1]
   }
}

// Update position pos (1-based) to value v
func (st *SegTree) Update(pos int, v int64) {
   idx := st.n + pos - 1
   st.tree[idx] = v
   for idx >>= 1; idx > 0; idx >>= 1 {
       st.tree[idx] = st.tree[2*idx] + st.tree[2*idx+1]
   }
}

// Query sum on [l..r] inclusive (1-based)
func (st *SegTree) Query(l, r int) int64 {
   l = l + st.n - 1
   r = r + st.n - 1
   var res int64
   for l <= r {
       if l&1 == 1 {
           res += st.tree[l]
           l++
       }
       if r&1 == 0 {
           res += st.tree[r]
           r--
       }
       l >>= 1
       r >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   a := make([]int64, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // diffs d[2..n]
   d := make([]int64, n+2)
   // posVals holds max(d[i],0)
   pos := make([]int64, n+2)
   for i := 2; i <= n; i++ {
       d[i] = a[i] - a[i-1]
       if d[i] > 0 {
           pos[i] = d[i]
       }
   }
   st := NewSegTree(n + 2)
   st.Build(pos)
   // a1
   var a1 int64 = a[1]
   // function to output answer
   printAns := func() {
       sumPos := st.Query(2, n)
       total := a1 + sumPos
       var ans int64
       if total >= 0 {
           ans = (total + 1) / 2
       } else {
           ans = total / 2
       }
       fmt.Fprintln(out, ans)
   }
   // initial answer
   printAns()
   var q int
   fmt.Fscan(in, &q)
   for i := 0; i < q; i++ {
       var l, r int
       var x int64
       fmt.Fscan(in, &l, &r, &x)
       if l == 1 {
           a1 += x
       }
       // update d[l]
       if l > 1 {
           // old d[l]
           old := d[l]
           d[l] = old + x
           var nv int64
           if d[l] > 0 {
               nv = d[l]
           }
           st.Update(l, nv)
       }
       // update d[r+1]
       if r+1 <= n {
           old := d[r+1]
           d[r+1] = old - x
           var nv int64
           if d[r+1] > 0 {
               nv = d[r+1]
           }
           st.Update(r+1, nv)
       }
       printAns()
   }
}
