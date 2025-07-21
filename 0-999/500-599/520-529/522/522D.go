package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = int(1e9)

func min(a, b int) int {
   if a < b {
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
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   // queriesByR[r] = list of (l, idx)
   type Query struct{ l, idx int }
   queriesByR := make([][]Query, n+1)
   ans := make([]int, m)
   for i := 0; i < m; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       queriesByR[r] = append(queriesByR[r], Query{l, i})
       ans[i] = -1
   }

   // iterative segment tree for range min
   size := 1
   for size < n+2 {
       size <<= 1
   }
   tree := make([]int, 2*size)
   for i := range tree {
       tree[i] = INF
   }

   // update position p to value v
   update := func(p, v int) {
       p += size - 1
       tree[p] = v
       for p > 1 {
           p >>= 1
           tree[p] = min(tree[2*p], tree[2*p+1])
       }
   }

   // query min on [l, r]
   query := func(l, r int) int {
       l += size - 1
       r += size - 1
       res := INF
       for l <= r {
           if l&1 == 1 {
               res = min(res, tree[l])
               l++
           }
           if r&1 == 0 {
               res = min(res, tree[r])
               r--
           }
           l >>= 1
           r >>= 1
       }
       return res
   }

   last := make(map[int]int)
   for i := 1; i <= n; i++ {
       x := a[i]
       if p, ok := last[x]; ok {
           dist := i - p
           update(p, dist)
       }
       last[x] = i
       // answer queries ending at i
       for _, q := range queriesByR[i] {
           res := query(q.l, i)
           if res >= INF {
               ans[q.idx] = -1
           } else {
               ans[q.idx] = res
           }
       }
   }

   for i := 0; i < m; i++ {
       fmt.Fprintln(writer, ans[i])
   }
}
