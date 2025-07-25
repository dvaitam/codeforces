package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int64
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   p := make([]int64, m)
   for i := int64(0); i < m; i++ {
       fmt.Fscan(reader, &p[i])
   }
   var removed int64 = 0
   var ops int64 = 0
   i := int64(0)
   for i < m {
       // determine current page index (0-based)
       curr := (p[i] - removed - 1) / k
       // count how many items in this page
       var cnt int64 = 0
       for i < m && (p[i]-removed-1)/k == curr {
           cnt++
           i++
       }
       // perform one operation
       removed += cnt
       ops++
   }
   fmt.Println(ops)
}
