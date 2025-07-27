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
   for i := 0; i < t; i++ {
       var a, b, c, d int
       fmt.Fscan(reader, &a, &b, &c, &d)
       // The minimal cutoff score is max(a + b, c + d)
       sum1 := a + b
       sum2 := c + d
       if sum2 > sum1 {
           sum1 = sum2
       }
       fmt.Fprintln(writer, sum1)
   }
}
