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
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       var s string
       fmt.Fscan(reader, &s)
       res := make([]byte, 0, n)
       for i := 0; i < n; i++ {
           idx := 2 * i
           if idx < len(s) {
               res = append(res, s[idx])
           }
       }
       fmt.Fprintln(writer, string(res))
   }
}
