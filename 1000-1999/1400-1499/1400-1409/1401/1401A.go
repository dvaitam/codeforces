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
   for t > 0 {
       t--
       var n, k int
       fmt.Fscan(reader, &n, &k)
       if n < k {
           fmt.Fprintln(writer, k-n)
       } else if (n-k)%2 == 0 {
           fmt.Fprintln(writer, 0)
       } else {
           fmt.Fprintln(writer, 1)
       }
   }
}
