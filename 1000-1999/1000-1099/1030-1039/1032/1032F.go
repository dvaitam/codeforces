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

   var s0 string
   fmt.Fscan(reader, &s0)
   // preallocate capacity for appended copies
   buf := make([]byte, len(s0), len(s0)+5000000)
   copy(buf, s0)

   var m int
   fmt.Fscan(reader, &m)
   for i := 0; i < m; i++ {
       var op int
       fmt.Fscan(reader, &op)
       if op == 1 {
           var l, r int
           fmt.Fscan(reader, &l, &r)
           var t string
           fmt.Fscan(reader, &t)
           buf = append(buf, t...)
       } else {
           var l, r, k int
           fmt.Fscan(reader, &l, &r, &k)
           idx := l + k - 2 // zero-based index
           writer.WriteByte(buf[idx])
           writer.WriteByte('\n')
       }
   }
}
