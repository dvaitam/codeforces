package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Shirt represents a t-shirt with cost and quality
type Shirt struct {
   c, q int64
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   shirts := make([]Shirt, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &shirts[i].c, &shirts[i].q)
   }
   // sort by descending quality, ascending cost
   sort.Slice(shirts, func(i, j int) bool {
       if shirts[i].q != shirts[j].q {
           return shirts[i].q > shirts[j].q
       }
       return shirts[i].c < shirts[j].c
   })
   // build prefix sums of costs
   prefix := make([]int64, n+1)
   for i := 0; i < n; i++ {
       prefix[i+1] = prefix[i] + shirts[i].c
   }
   var k int
   fmt.Fscan(reader, &k)
   ans := make([]int, k)
   for i := 0; i < k; i++ {
       var b int64
       fmt.Fscan(reader, &b)
       // binary search maximum t such that prefix[t] <= b
       lo, hi := 0, n
       for lo < hi {
           mid := (lo + hi + 1) >> 1
           if prefix[mid] <= b {
               lo = mid
           } else {
               hi = mid - 1
           }
       }
       ans[i] = lo
   }
   // output answers
   for i, v := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}
