package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n, m int
       fmt.Fscan(in, &n, &m)
       // Print even-indexed rows first
       for i := 2; i <= n; i += 2 {
           base := (i-1) * m
           for j := 1; j <= m; j++ {
               fmt.Fprint(out, base+j)
               if j < m {
                   out.WriteByte(' ')
               }
           }
           out.WriteByte('\n')
       }
       // Then odd-indexed rows
       for i := 1; i <= n; i += 2 {
           base := (i-1) * m
           for j := 1; j <= m; j++ {
               fmt.Fprint(out, base+j)
               if j < m {
                   out.WriteByte(' ')
               }
           }
           out.WriteByte('\n')
       }
   }
}
