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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tt := 0; tt < t; tt++ {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int64, n)
       var total int64
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
           total += a[i]
       }
       // find maximum number of equal-sum segments
       maxSeg := int64(1)
       var prefixSum int64
       for i := 0; i < n; i++ {
           prefixSum += a[i]
           segSum := prefixSum
           if total%segSum != 0 {
               continue
           }
           need := total / segSum
           var cur int64
           var cnt int64
           ok := true
           for j := 0; j < n; j++ {
               cur += a[j]
               if cur == segSum {
                   cnt++
                   cur = 0
               } else if cur > segSum {
                   ok = false
                   break
               }
           }
           if ok && cnt == need && need > maxSeg {
               maxSeg = need
           }
       }
       // operations = n - number of segments
       ops := n - int(maxSeg)
       fmt.Fprintln(writer, ops)
   }
}
