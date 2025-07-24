package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var c int64
   if _, err := fmt.Fscan(reader, &n, &c); err != nil {
       return
   }
   var last int64
   var t int64
   count := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &t)
       if i == 0 || t-last <= c {
           count++
       } else {
           count = 1
       }
       last = t
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, count)
}
