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
   fmt.Fscan(reader, &t)
   for i := 0; i < t; i++ {
       var x1, y1, z1 int64
       var x2, y2, z2 int64
       fmt.Fscan(reader, &x1, &y1, &z1)
       fmt.Fscan(reader, &x2, &y2, &z2)
       // maximize 2*min(z1, y2) positive pairs
       pos := minInt64(z1, y2)
       // unavoidable negative pairs: a_i=1 with b_i=2
       // safe spots for a's 1's: b's zeros and ones
       safe := x2 + y2
       negPairs := y1 - safe
       if negPairs < 0 {
           negPairs = 0
       }
       result := 2*pos - 2*negPairs
       fmt.Fprintln(writer, result)
   }
}

func minInt64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}
