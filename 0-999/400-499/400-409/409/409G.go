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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var x, y float64
   var ySum float64
   for i := 0; i < n; i++ {
       if _, err := fmt.Fscan(reader, &x, &y); err != nil {
           return
       }
       ySum += y
   }
   ans := 5 + ySum/float64(n)
   fmt.Fprintf(writer, "%.10f\n", ans)
}
