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
   w := make([]int64, n)
   h := make([]int64, n)
   var sumW int64
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &w[i], &h[i])
       sumW += w[i]
   }

   // find the largest and second largest heights
   var maxH, secondH int64
   var idxMax int
   for i, hi := range h {
       if hi > maxH {
           secondH = maxH
           maxH = hi
           idxMax = i
       } else if hi > secondH {
           secondH = hi
       }
   }

   // compute result for each friend
   for i := 0; i < n; i++ {
       curH := maxH
       if i == idxMax {
           curH = secondH
       }
       area := (sumW - w[i]) * curH
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, area)
   }
}
