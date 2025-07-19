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

   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for T > 0 {
       T--
       var n int
       fmt.Fscan(reader, &n)
       xor := 0
       last := 0
       for i := 0; i < n; i++ {
           var v int
           fmt.Fscan(reader, &v)
           xor ^= v
           last = v
       }
       if n%2 == 1 {
           fmt.Fprintln(writer, xor)
       } else {
           if xor == 0 {
               fmt.Fprintln(writer, last)
           } else {
               fmt.Fprintln(writer, -1)
           }
       }
   }
}
