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
       var w, h, n int64
       fmt.Fscan(reader, &w, &h, &n)
       var cnt int64 = 1
       for w%2 == 0 {
           cnt <<= 1
           w >>= 1
       }
       for h%2 == 0 {
           cnt <<= 1
           h >>= 1
       }
       if cnt >= n {
           fmt.Fprintln(writer, "YES")
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}
