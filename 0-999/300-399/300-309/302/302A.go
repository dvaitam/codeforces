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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   cnt1 := 0
   // total -1 count = n - cnt1
   for i := 0; i < n; i++ {
       var v int
       fmt.Fscan(reader, &v)
       if v == 1 {
           cnt1++
       }
   }
   cntNeg1 := n - cnt1
   // process queries
   for i := 0; i < m; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       length := r - l + 1
       if length%2 != 0 {
           fmt.Fprint(writer, "0")
       } else {
           need := length / 2
           if cnt1 >= need && cntNeg1 >= need {
               fmt.Fprint(writer, "1")
           } else {
               fmt.Fprint(writer, "0")
           }
       }
       if i < m-1 {
           fmt.Fprint(writer, "\n")
       }
   }
}
