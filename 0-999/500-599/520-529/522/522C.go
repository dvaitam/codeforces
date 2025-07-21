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
   for ; t > 0; t-- {
       var m, k int
       fmt.Fscan(reader, &m, &k)
       a := make([]int, k+1)
       for i := 1; i <= k; i++ {
           fmt.Fscan(reader, &a[i])
       }
       // count known consumptions and unknown slots
       cnt := make([]int, k+1)
       unknown := 0
       for i := 0; i < m-1; i++ {
           var ti, ri int
           fmt.Fscan(reader, &ti, &ri)
           if ti > 0 {
               if ti <= k {
                   cnt[ti]++
               }
           } else {
               unknown++
           }
       }
       // output result
       for i := 1; i <= k; i++ {
           if cnt[i]+unknown >= a[i] {
               writer.WriteByte('Y')
           } else {
               writer.WriteByte('N')
           }
       }
       writer.WriteByte('\n')
   }
}
