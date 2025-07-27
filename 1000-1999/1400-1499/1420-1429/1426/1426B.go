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
       var n, m int
       fmt.Fscan(reader, &n, &m)
       hasSym := false
       for j := 0; j < n; j++ {
           var a, b, c, d int
           fmt.Fscan(reader, &a, &b)
           fmt.Fscan(reader, &c, &d)
           if b == c {
               hasSym = true
           }
       }
       if m%2 == 1 {
           fmt.Fprintln(writer, "NO")
       } else {
           if hasSym {
               fmt.Fprintln(writer, "YES")
           } else {
               fmt.Fprintln(writer, "NO")
           }
       }
   }
}
