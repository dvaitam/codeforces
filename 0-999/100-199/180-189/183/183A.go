package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var cntUL, cntUR, cntDL, cntDR int64
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       switch s {
       case "UL":
           cntUL++
       case "UR":
           cntUR++
       case "DL":
           cntDL++
       case "DR":
           cntDR++
       // "ULDR" or other strings do not affect these counters
       }
   }
   // Calculate free moves counts for u and v dimensions
   total := int64(n)
   freeU := total - cntUR - cntDL
   freeV := total - cntUL - cntDR
   // Number of possible sums is free + 1
   ans := (freeU + 1) * (freeV + 1)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
}
