package main

import (
   "bufio"
   "fmt"
   "os"
)

type Item struct {
   c, h int
}

var (
   tree [][]Item
   qs   [][][2]int // per time: list of (query index, budget)
   ans  []int
   Bmax int
)

func addItem(idx, l, r, ql, qr int, it Item) {
   if ql > r || qr < l {
       return
   }
   if ql <= l && r <= qr {
       tree[idx] = append(tree[idx], it)
       return
   }
   m := (l + r) >> 1
   addItem(idx<<1, l, m, ql, qr, it)
   addItem(idx<<1|1, m+1, r, ql, qr, it)
}

func dfs(idx, l, r int, dp []int32) {
   // apply items at this node
   var dpCur []int32 = dp
   if len(tree[idx]) > 0 {
       dpNew := make([]int32, Bmax+1)
       copy(dpNew, dpCur)
       for _, it := range tree[idx] {
           ci, hi := it.c, it.h
           for b := Bmax; b >= ci; b-- {
               v := dpNew[b-ci] + int32(hi)
               if v > dpNew[b] {
                   dpNew[b] = v
               }
           }
       }
       dpCur = dpNew
   }
   if l == r {
       for _, qb := range qs[l] {
           qi, b := qb[0], qb[1]
           if b > Bmax {
               b = Bmax
           }
           ans[qi] = int(dpCur[b])
       }
       return
   }
   m := (l + r) >> 1
   dfs(idx<<1, l, m, dpCur)
   dfs(idx<<1|1, m+1, r, dpCur)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, p int
   fmt.Fscan(in, &n, &p)
   items := make([][3]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &items[i][0], &items[i][1], &items[i][2])
   }
   var q int
   fmt.Fscan(in, &q)
   ans = make([]int, q)
   times := make([]int, q)
   budgets := make([]int, q)
   maxT := 0
   Bmax = 0
   qs = make([][][2]int, 1) // dummy at 0
   // read queries first to know max time and Bmax
   for i := 0; i < q; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       times[i], budgets[i] = a, b
       if a > maxT {
           maxT = a
       }
       if b > Bmax {
           Bmax = b
       }
       qs = append(qs, [][2]int{})
   }
   // ensure qs size
   if len(qs) <= maxT {
       more := make([][][2]int, maxT+1-len(qs))
       qs = append(qs, more...)
   }
   for i := 0; i < q; i++ {
       t, b := times[i], budgets[i]
       qs[t] = append(qs[t], [2]int{i, b})
   }
   // build segment tree
   size := 1
   for size <= maxT {
       size <<= 1
   }
   tree = make([][]Item, size<<1)
   // add items
   for _, it := range items {
       c, h, ti := it[0], it[1], it[2]
       l := ti
       r := ti + p - 1
       if r < 1 || l > maxT {
           continue
       }
       if l < 1 {
           l = 1
       }
       if r > maxT {
           r = maxT
       }
       addItem(1, 1, maxT, l, r, Item{c: c, h: h})
   }
   // initial dp
   dp0 := make([]int32, Bmax+1)
   // dfs
   dfs(1, 1, maxT, dp0)
   // output answers
   for i := 0; i < q; i++ {
       fmt.Fprintln(out, ans[i])
   }
}
