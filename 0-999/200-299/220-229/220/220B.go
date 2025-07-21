package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "math"
)

type Query struct {
   l, r, idx, blk int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if x <= n {
           a[i] = x
       } else {
           a[i] = 0
       }
   }
   queries := make([]Query, m)
   for i := 0; i < m; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       l-- // zero-based
       r--
       queries[i] = Query{l: l, r: r, idx: i}
   }
   // Mo's algorithm
   sz := int(math.Sqrt(float64(n)))
   if sz == 0 {
       sz = 1
   }
   for i := range queries {
       queries[i].blk = queries[i].l / sz
   }
   sort.Slice(queries, func(i, j int) bool {
       bi, bj := queries[i].blk, queries[j].blk
       if bi != bj {
           return bi < bj
       }
       return queries[i].r < queries[j].r
   })
   freq := make([]int, n+1)
   ansArr := make([]int, m)
   curL, curR, curAns := 0, -1, 0
   add := func(pos int) {
       x := a[pos]
       if x == 0 {
           return
       }
       if freq[x] == x {
           curAns--
       }
       freq[x]++
       if freq[x] == x {
           curAns++
       }
   }
   remove := func(pos int) {
       x := a[pos]
       if x == 0 {
           return
       }
       if freq[x] == x {
           curAns--
       }
       freq[x]--
       if freq[x] == x {
           curAns++
       }
   }
   for _, q := range queries {
       L, R := q.l, q.r
       for curL > L {
           curL--
           add(curL)
       }
       for curR < R {
           curR++
           add(curR)
       }
       for curL < L {
           remove(curL)
           curL++
       }
       for curR > R {
           remove(curR)
           curR--
       }
       ansArr[q.idx] = curAns
   }
   for i := 0; i < m; i++ {
       fmt.Fprintln(writer, ansArr[i])
   }
}
