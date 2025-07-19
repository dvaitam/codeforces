package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       pos1, pos2, posn := -1, -1, -1
       for i := 1; i <= n; i++ {
           var x int
           fmt.Fscan(reader, &x)
           if x == 1 {
               pos1 = i
           } else if x == 2 {
               pos2 = i
           } else if x == n {
               posn = i
           }
       }
       var i, j int
       low := min(pos1, pos2)
       high := max(pos1, pos2)
       if low < posn && posn < high {
           i, j = pos1, pos2
       } else if posn < low {
           i, j = posn, low
       } else {
           i, j = high, posn
       }
       fmt.Fprintln(writer, i, j)
   }
}
