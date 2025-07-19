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
   for t > 0 {
       t--
       var n int
       var d, h int64
       fmt.Fscan(reader, &n, &d, &h)
       y := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &y[i])
       }
       area := float64(d*h) / 2.0
       totArea := float64(n) * area
       for i := 1; i < n; i++ {
           diff := float64(y[i] - y[i-1])
           if diff < float64(h) {
               h2 := float64(h) - diff
               d2 := h2 * float64(d) / float64(h)
               ext := d2 * h2 / 2.0
               totArea -= ext
           }
       }
       fmt.Fprintf(writer, "%.7f\n", totArea)
   }
}
