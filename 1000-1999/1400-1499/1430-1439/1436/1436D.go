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
   fmt.Fscan(reader, &n)
   // build tree
   children := make([][]int, n+1)
   for i := 2; i <= n; i++ {
       var p int
       fmt.Fscan(reader, &p)
       children[p] = append(children[p], i)
   }
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // leaves count and subtree sum
   leaves := make([]int64, n+1)
   sum := make([]int64, n+1)
   var ans int64
   // process nodes in reverse order: children have higher indices
   for u := n; u >= 1; u-- {
       if len(children[u]) == 0 {
           leaves[u] = 1
           sum[u] = a[u]
       } else {
           var lcnt int64
           var s int64
           for _, v := range children[u] {
               lcnt += leaves[v]
               s += sum[v]
           }
           leaves[u] = lcnt
           sum[u] = s + a[u]
       }
       // compute minimal max load in this subtree
       // max of current ans and ceil(sum/leaves)
       // avoid division by zero: leaves[u] >=1
       x := (sum[u] + leaves[u] - 1) / leaves[u]
       if x > ans {
           ans = x
       }
   }
   fmt.Fprintln(writer, ans)
}
