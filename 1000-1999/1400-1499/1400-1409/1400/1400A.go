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
       var n int
       var s string
       fmt.Fscan(reader, &n)
       fmt.Fscan(reader, &s)
       cnt0 := 0
       for i := 0; i < len(s); i++ {
           if s[i] == '0' {
               cnt0++
           }
       }
       if cnt0 >= n {
           for i := 0; i < n; i++ {
               writer.WriteByte('0')
           }
       } else {
           for i := 0; i < n; i++ {
               writer.WriteByte('1')
           }
       }
       writer.WriteByte('\n')
   }
}
