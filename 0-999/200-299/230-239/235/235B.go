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
   var m, r float64
   for i := 0; i < n; i++ {
       var p float64
       if _, err := fmt.Fscan(reader, &p); err != nil {
           return
       }
       r += p * (2*m + 1)
       m = (m + 1) * p
   }
   fmt.Fprintf(writer, "%.10f\n", r)
}
