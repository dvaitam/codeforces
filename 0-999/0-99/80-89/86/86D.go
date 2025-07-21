package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

type Query struct {
   l, r, idx, block int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, t int
   if _, err := fmt.Fscan(reader, &n, &t); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   queries := make([]Query, t)
   for i := 0; i < t; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       l-- // zero-index
       r-- // zero-index
       queries[i] = Query{l: l, r: r, idx: i}
   }
   // Mo's algorithm block size
   blockSize := int(math.Sqrt(float64(n)))
   for i := range queries {
       queries[i].block = queries[i].l / blockSize
   }
   sort.Slice(queries, func(i, j int) bool {
       qi, qj := queries[i], queries[j]
       if qi.block != qj.block {
           return qi.block < qj.block
       }
       if qi.block%2 == 0 {
           return qi.r < qj.r
       }
       return qi.r > qj.r
   })
   // prepare for processing
   counts := make([]int, 1000001)
   ansArr := make([]int64, t)
   var currL, currR int = 0, -1
   var currAns int64 = 0
   add := func(pos int) {
       v := a[pos]
       c := counts[v]
       currAns += int64(2*c+1) * int64(v)
       counts[v] = c + 1
   }
   remove := func(pos int) {
       v := a[pos]
       c := counts[v]
       currAns += int64(-2*c+1) * int64(v)
       counts[v] = c - 1
   }
   // process queries
   for _, q := range queries {
       L, R := q.l, q.r
       for currL > L {
           currL--
           add(currL)
       }
       for currR < R {
           currR++
           add(currR)
       }
       for currL < L {
           remove(currL)
           currL++
       }
       for currR > R {
           remove(currR)
           currR--
       }
       ansArr[q.idx] = currAns
   }
   // output answers
   for i := 0; i < t; i++ {
       fmt.Fprintln(writer, ansArr[i])
   }
}
