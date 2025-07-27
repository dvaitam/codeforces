package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for i := 0; i < t; i++ {
       var a, b, x, y, n int64
       fmt.Fscan(reader, &a, &b, &x, &y, &n)
       // Option 1: decrease a first, then b
       decA := min(n, a-x)
       a1 := a - decA
       rem1 := n - decA
       decB := min(rem1, b-y)
       b1 := b - decB

       // Option 2: decrease b first, then a
       decB2 := min(n, b-y)
       b2 := b - decB2
       rem2 := n - decB2
       decA2 := min(rem2, a-x)
       a2 := a - decA2

       prod1 := a1 * b1
       prod2 := a2 * b2
       if prod1 < prod2 {
           fmt.Fprintln(writer, prod1)
       } else {
           fmt.Fprintln(writer, prod2)
       }
   }
}
