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
      var l, r int64
      fmt.Fscan(reader, &l, &r)
      var ans int64
      for p := int64(1); p <= r; p *= 10 {
         ans += r/p - l/p
      }
      fmt.Fprintln(writer, ans)
   }
}
