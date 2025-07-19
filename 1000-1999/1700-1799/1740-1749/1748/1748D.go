package main

import (
   "bufio"
   "fmt"
   "math/bits"
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
   for i := 0; i < t; i++ {
       var a, b, d uint64
       fmt.Fscan(reader, &a, &b, &d)
       if bits.TrailingZeros64(a|b) < bits.TrailingZeros64(d) {
           fmt.Fprintln(writer, -1)
           continue
       }
       k := bits.TrailingZeros64(d)
       var x uint64
       for j := k; j < 30; j++ {
           if ((^x) >> j & 1) == 1 {
               x += d << (j - k)
           }
       }
       fmt.Fprintln(writer, x)
   }
}
