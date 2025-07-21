package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, c int64
   if _, err := fmt.Fscan(reader, &a, &c); err != nil {
       return
   }
   var b int64 = 0
   var base int64 = 1
   for a > 0 || c > 0 {
       da := a % 3
       dc := c % 3
       // b digit so that (da + db) % 3 == dc
       db := (dc - da + 3) % 3
       b += db * base
       base *= 3
       a /= 3
       c /= 3
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, b)
}
