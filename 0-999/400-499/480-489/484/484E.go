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
   stack := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       for len(stack) > 0 && h[stack[len(stack)-1]] >= h[i] {
           stack = stack[:len(stack)-1]
       }
       if len(stack) == 0 {
           L[i] = 0
       } else {
           L[i] = stack[len(stack)-1]
       }
       stack = append(stack, i)
   }
   // compute R[i]
   R := make([]int, n+1)
   stack = stack[:0]
   for i := n; i >= 1; i-- {
       for len(stack) > 0 && h[stack[len(stack)-1]] >= h[i] {
           stack = stack[:len(stack)-1]
       }
       if len(stack) == 0 {
           R[i] = n + 1
       } else {
           R[i] = stack[len(stack)-1]
       }
       stack = append(stack, i)
   }
   // build bars with width
   type BarW struct{ a, b, w, h int }
   barw := make([]BarW, n)
   for i := 1; i <= n; i++ {
       barw[i-1] = BarW{L[i], R[i], R[i] - L[i] - 1, h[i]}
   }
   // sort bars by width descending
   sort.Slice(barw, func(i, j int) bool { return barw[i].w > barw[j].w })

   // read queries
   var m int
   fmt.Fscan(in, &m)
   type QueryW struct{ w, t1, t2, idx int }
   qs2 := make([]QueryW, m)
   res := make([]int, m)
   for i := 0; i < m; i++ {
       var l, r, w int
       fmt.Fscan(in, &l, &r, &w)
       qs2[i] = QueryW{w, r - w, l + w, i}
   }
   // sort queries by width descending
   sort.Slice(qs2, func(i, j int) bool { return qs2[i].w > qs2[j].w })

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
   for _, bw := range barw {
       collect(1, 0, n, bw.a, bw.b)
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
       resL := queryTree(v<<1, l, m2, ql, qr, T2)
       resR := queryTree(v<<1|1, m2+1, r, ql, qr, T2)
       if resL > resR {
           return resL
       }
       return resR
   }

   // process queries and bars
   bi := 0
   for _, q := range qs2 {
       // insert bars with width >= q.w
       for bi < n && barw[bi].w >= q.w {
           update(1, 0, n, barw[bi].a, barw[bi].b, barw[bi].h)
           bi++
       }
       // query a in [0..q.t1], b >= q.t2
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
