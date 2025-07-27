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
   counts := make([]int, 20)
   for i := 0; i < n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       for b := 0; b < 20; b++ {
           if (a>>b)&1 == 1 {
               counts[b]++
           }
       }
   }
   var result uint64 = 0
   // Build numbers greedily
   for i := 0; i < n; i++ {
       var x int
       for b := 0; b < 20; b++ {
           if counts[b] > 0 {
               x |= 1 << b
               counts[b]--
           }
       }
       result += uint64(x) * uint64(x)
   }
   fmt.Fprint(writer, result)
}
