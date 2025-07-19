package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Query holds a request with start position, type (negated), and original index
type Query struct {
   a   int
   b   int
   idx int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var N, M int
   if _, err := fmt.Fscan(reader, &N, &M); err != nil {
       return
   }
   qs := make([]Query, M)
   for i := 0; i < M; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       qs[i] = Query{a: a, b: -b, idx: i}
   }

   sort.Slice(qs, func(i, j int) bool {
       if qs[i].a != qs[j].a {
           return qs[i].a < qs[j].a
       }
       if qs[i].b != qs[j].b {
           return qs[i].b < qs[j].b
       }
       return qs[i].idx < qs[j].idx
   })

   res1 := make([]int, M)
   res2 := make([]int, M)
   con, ua, ub := 0, 0, 2

   for _, q := range qs {
       if q.b == 0 {
           // type 0 request
           if ub > con {
               fmt.Fprintln(writer, -1)
               return
           }
           r1, r2 := ua, ub
           ua++
           if ua+1 == ub {
               ua = 0
               ub++
           }
           res1[q.idx] = r1 + 1
           res2[q.idx] = r2 + 1
       } else {
           // type 1 request
           r1, r2 := con, con+1
           con++
           res1[q.idx] = r1 + 1
           res2[q.idx] = r2 + 1
       }
   }

   for i := 0; i < M; i++ {
       fmt.Fprintln(writer, res1[i], res2[i])
   }
}
