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
   for ; t > 0; t-- {
       var n int
       var s string
       fmt.Fscan(reader, &n, &s)
       ok := true
       for j := 0; j < n/2; j++ {
           diff := s[j] - s[n-1-j]
           if diff < 0 {
               diff = -diff
           }
           if diff != 0 && diff != 2 {
               ok = false
               break
           }
       }
       if ok {
           fmt.Fprintln(writer, "YES")
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}
