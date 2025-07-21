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
   leftBound, rightBound := 1, n
   for i := 0; i < m; i++ {
       // hint format: To the left of i  OR To the right of i
       var w1, w2, dir, w4 string
       var idx int
       fmt.Fscan(reader, &w1, &w2, &dir, &w4, &idx)
       if dir == "left" {
           // hidden in [1, idx-1]
           if idx-1 < rightBound {
               rightBound = idx - 1
           }
       } else if dir == "right" {
           // hidden in [idx+1, n]
           if idx+1 > leftBound {
               leftBound = idx + 1
           }
       }
   }
   if leftBound > rightBound {
       fmt.Fprintln(writer, -1)
   } else {
       fmt.Fprintln(writer, rightBound-leftBound+1)
   }
}
