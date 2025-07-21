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
   t := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &t[i])
   }

   prefixNonNeg := make([]int, n+1)
   for i := 1; i <= n; i++ {
       prefixNonNeg[i] = prefixNonNeg[i-1]
       if t[i-1] >= 0 {
           prefixNonNeg[i]++
       }
   }

   suffixNonPos := make([]int, n+2)
   for i := n; i >= 1; i-- {
       suffixNonPos[i] = suffixNonPos[i+1]
       if t[i-1] <= 0 {
           suffixNonPos[i]++
       }
   }

   ans := n
   for k := 1; k <= n-1; k++ {
       cost := prefixNonNeg[k] + suffixNonPos[k+1]
       if cost < ans {
           ans = cost
       }
   }
   fmt.Fprintln(writer, ans)
}
