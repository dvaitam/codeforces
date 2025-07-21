package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n+2)
   for i := 1; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // lft[i]: max usable crossings in prefix [1..i] to return to endpoint i
   lft := make([]int64, n+1)
   for i := 1; i < n; i++ {
       if lft[i-1]&1 == 1 {
           lft[i] = lft[i-1] + a[i] - 1
       } else {
           lft[i] = lft[i-1] + a[i]
       }
   }
   // rgt[i]: max usable crossings in suffix [i..n-1] to return to endpoint i
   rgt := make([]int64, n+2)
   for i := n - 1; i >= 1; i-- {
       if rgt[i+1]&1 == 1 {
           rgt[i] = rgt[i+1] + a[i] - 1
       } else {
           rgt[i] = rgt[i+1] + a[i]
       }
   }
   var ans int64
   // consider each bridge i as last collapse
   for i := 1; i < n; i++ {
       total := lft[i-1] + a[i] + rgt[i+1]
       if total > ans {
           ans = total
       }
   }
   fmt.Fprintln(writer, ans)
}
