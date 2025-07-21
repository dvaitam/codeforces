package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, k int
   var b int64
   fmt.Fscan(reader, &n, &k)
   fmt.Fscan(reader, &b)
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // can secure square x?
   can := func(x int) bool {
       // build drain candidates: all indices except x and except worst free
       drains := make([]int64, 0, n)
       if x < n {
           // worst free is n
           for i := 1; i <= n; i++ {
               if i == x || i == n {
                   continue
               }
               drains = append(drains, a[i])
           }
       } else {
           // x == n, worst free is x, cannot pick x
           for i := 1; i < n; i++ {
               if i == x {
                   continue
               }
               drains = append(drains, a[i])
           }
       }
       // need to pick k-1 largest from drains
       need := k - 1
       if need <= 0 {
           // no drain, check directly
           return b < a[x]
       }
       if len(drains) == 0 {
           // nothing to drain
           return b < a[x]
       }
       // sort descending
       sort.Slice(drains, func(i, j int) bool { return drains[i] > drains[j] })
       var s int64
       for i := 0; i < need && i < len(drains); i++ {
           s += drains[i]
           if s > b {
               break
           }
       }
       // after drains, remaining budget r = b - s
       // check if r < a[x]
       return s > b - a[x]
   }

   // binary search smallest x in [1,n]
   l, r := 1, n
   ans := n
   for l <= r {
       mid := (l + r) / 2
       if can(mid) {
           ans = mid
           r = mid - 1
       } else {
           l = mid + 1
       }
   }
   fmt.Fprintln(writer, ans)
}
