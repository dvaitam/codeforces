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
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       sum := 0.0
       maxn := -1e18
       for i := 0; i < n; i++ {
           var x float64
           fmt.Fscan(reader, &x)
           if x > maxn {
               maxn = x
           }
           sum += x
       }
       sum -= maxn
       result := sum/float64(n-1) + maxn
       fmt.Fprintf(writer, "%.8f\n", result)
   }
}
