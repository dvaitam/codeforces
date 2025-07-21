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

   var n, m int
   fmt.Fscan(reader, &n, &m)
   prev := byte(0)
   ok := true
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       if len(s) != m {
           ok = false
           break
       }
       c := s[0]
       for j := 1; j < m; j++ {
           if s[j] != c {
               ok = false
               break
           }
       }
       if !ok {
           break
       }
       if i > 0 && c == prev {
           ok = false
           break
       }
       prev = c
   }
   if ok {
       fmt.Fprintln(writer, "YES")
   } else {
       fmt.Fprintln(writer, "NO")
   }
}
