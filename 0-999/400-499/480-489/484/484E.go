package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// BIT for range prefix max query and point update
type BIT struct {
   n int
   t []int
}

func NewBIT(n int) *BIT {
   return &BIT{n: n, t: make([]int, n+1)}
}

// Update position i (1-indexed) with value v: t[i] = max(t[i], v)
func (b *BIT) Update(i, v int) {
   for ; i <= b.n; i += i & -i {
       if b.t[i] < v {
           b.t[i] = v
       }
   }
}

// Query max on prefix [1..i]
func (b *BIT) Query(i int) int {
   res := 0
   for ; i > 0; i -= i & -i {
       if b.t[i] > res {
           res = b.t[i]
       }
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   h := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &h[i])
   }
   // compute L[i]
   L := make([]int, n+1)
   stk := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       for len(stk) > 0 && h[stk[len(stk)-1]] >= h[i] {
           stk = stk[:len(stk)-1]
       }
       if len(stk) == 0 {
           L[i] = 0
       } else {
           L[i] = stk[len(stk)-1]
       }
       stk = append(stk, i)
   }
   // compute R[i]
   R := make([]int, n+1)
   stk = stk[:0]
   for i := n; i >= 1; i-- {
       for len(stk) > 0 && h[stk[len(stk)-1]] >= h[i] {
           stk = stk[:len(stk)-1]
       }
       if len(stk) == 0 {
           R[i] = n + 1
       } else {
           R[i] = stk[len(stk)-1]
       }
       stk = append(stk, i)
   }
   // prepare bars with width
   type Bar struct{ a, b, w, h int }
   bars := make([]Bar, n)
   for i := 1; i <= n; i++ {
       bars[i-1] = Bar{L[i], R[i], R[i] - L[i] - 1, h[i]}
   }
   // sort bars by width descending
   sort.Slice(bars, func(i, j int) bool { return bars[i].w > bars[j].w })

   // read queries
   var m int
   fmt.Fscan(in, &m)
   type Query struct{ w, t1, t2, idx int }
   qs := make([]Query, m)
   res := make([]int, m)
   for i := 0; i < m; i++ {
       var l, r, w int
       fmt.Fscan(in, &l, &r, &w)
       qs[i] = Query{w, r - w, l + w, i}
   }
   // sort queries by width descending
   sort.Slice(qs, func(i, j int) bool { return qs[i].w > qs[j].w })

   // build segment tree on a in [0..n]
   sizeA := n + 1
   treeSize := 4 * sizeA
   type Node struct{ bs []int; bit *BIT }
   tree := make([]Node, treeSize)
   var collect func(v, l, r, ai, bi int)
   collect = func(v, l, r, ai, bi int) {
       tree[v].bs = append(tree[v].bs, bi)
       if l == r {
           return
       }
       m2 := (l + r) >> 1
       if ai <= m2 {
           collect(v<<1, l, m2, ai, bi)
       } else {
           collect(v<<1|1, m2+1, r, ai, bi)
       }
   }
   for _, b := range bars {
       collect(1, 0, n, b.a, b.b)
   }
   var initNode func(v, l, r int)
   initNode = func(v, l, r int) {
       if len(tree[v].bs) > 0 {
           bs := tree[v].bs
           sort.Ints(bs)
           u := 0
           for i := 1; i < len(bs); i++ {
               if bs[i] != bs[u] {
                   u++
                   bs[u] = bs[i]
               }
           }
           tree[v].bs = bs[:u+1]
           tree[v].bit = NewBIT(len(tree[v].bs))
       }
       if l == r {
           return
       }
       m2 := (l + r) >> 1
       initNode(v<<1, l, m2)
       initNode(v<<1|1, m2+1, r)
   }
   initNode(1, 0, n)

   var update func(v, l, r, ai, bi, hi int)
   update = func(v, l, r, ai, bi, hi int) {
       bs := tree[v].bs
       if len(bs) > 0 {
           pos := sort.Search(len(bs), func(i int) bool { return bs[i] >= bi })
           if pos < len(bs) {
               id2 := len(bs) - pos
               tree[v].bit.Update(id2, hi)
           }
       }
       if l == r {
           return
       }
       m2 := (l + r) >> 1
       if ai <= m2 {
           update(v<<1, l, m2, ai, bi, hi)
       } else {
           update(v<<1|1, m2+1, r, ai, bi, hi)
       }
   }

   var queryTree func(v, l, r, ql, qr, T2 int) int
   queryTree = func(v, l, r, ql, qr, T2 int) int {
       if qr < l || r < ql {
           return 0
       }
       if ql <= l && r <= qr {
           bs := tree[v].bs
           if len(bs) == 0 {
               return 0
           }
           pos := sort.Search(len(bs), func(i int) bool { return bs[i] >= T2 })
           if pos == len(bs) {
               return 0
           }
           id2 := len(bs) - pos
           return tree[v].bit.Query(id2)
       }
       m2 := (l + r) >> 1
       v1 := queryTree(v<<1, l, m2, ql, qr, T2)
       v2 := queryTree(v<<1|1, m2+1, r, ql, qr, T2)
       if v1 > v2 {
           return v1
       }
       return v2
   }

   // process queries and bars
   bi2 := 0
   for _, q := range qs {
       for bi2 < n && bars[bi2].w >= q.w {
           update(1, 0, n, bars[bi2].a, bars[bi2].b, bars[bi2].h)
           bi2++
       }
       if q.t1 >= 0 {
           res[q.idx] = queryTree(1, 0, n, 0, q.t1, q.t2)
       } else {
           res[q.idx] = 0
       }
   }
   for i := 0; i < m; i++ {
       fmt.Fprintln(out, res[i])
   }
}
