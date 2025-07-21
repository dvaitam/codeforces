package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

const INF = int(1e9)

// Segment tree for range minimum query and point update
type SegTree struct {
   n    int
   tree []int
}

func NewSegTree(n int) *SegTree {
   size := 1
   for size < n {
       size <<= 1
   }
   tree := make([]int, 2*size)
   for i := range tree {
       tree[i] = INF
   }
   return &SegTree{n: size, tree: tree}
}

func (st *SegTree) Update(pos, val int) {
   i := pos + st.n
   if st.tree[i] <= val {
       return
   }
   st.tree[i] = val
   for i >>= 1; i > 0; i >>= 1 {
       st.tree[i] = min(st.tree[2*i], st.tree[2*i+1])
   }
}

func (st *SegTree) Query(l, r int) int {
   // inclusive l,r
   l += st.n
   r += st.n
   res := INF
   for l <= r {
       if l&1 == 1 {
           res = min(res, st.tree[l])
           l++
       }
       if r&1 == 0 {
           res = min(res, st.tree[r])
           r--
       }
       l >>= 1
       r >>= 1
   }
   return res
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   x := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x[i])
   }
   var a, b int
   fmt.Fscan(reader, &a, &b)
   L := a - b
   st := NewSegTree(L + 1)
   // initialize dp[d] = INF, except dp[L] = 0
   st.Update(L, 0)
   // dp by point reads from seg
   for d := L; d >= 1; d-- {
       v := d // index
       cur := st.Query(d, d)
       // move -1
       prev := st.Query(d-1, d-1)
       if cur+1 < prev {
           st.Update(d-1, cur+1)
       }
       // mod moves
       realV := b + d
       for _, xi := range x {
           m := realV - realV%xi
           if m < b {
               continue
           }
           dm := m - b
           old := st.Query(dm, dm)
           if cur+1 < old {
               st.Update(dm, cur+1)
           }
       }
   }
   ans := st.Query(0, 0)
   fmt.Println(ans)
}
