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
   // counts for t from -10 to 10, index offset by +10
   var cnt [21]int64
   for i := 0; i < n; i++ {
       var t int
       fmt.Fscan(reader, &t)
       if t < -10 || t > 10 {
           // skip invalid, though per constraints shouldn't happen
           continue
       }
       cnt[t+10]++
   }
   var res int64
   // pairs of zeros
   z := cnt[0+10]
   if z > 1 {
       res += z * (z - 1) / 2
   }
   // pairs for x and -x, x from 1 to 10
   for x := 1; x <= 10; x++ {
       res += cnt[x+10] * cnt[-x+10]
   }
   fmt.Fprintln(writer, res)
}
