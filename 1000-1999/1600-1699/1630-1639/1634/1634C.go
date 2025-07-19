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

   var t, n, k int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       fmt.Fscan(reader, &n, &k)
       if k == 1 {
           fmt.Fprintln(writer, "YES")
           for i := 1; i <= n; i++ {
               fmt.Fprintln(writer, i)
           }
       } else if n%2 == 1 {
           fmt.Fprintln(writer, "NO")
       } else {
           fmt.Fprintln(writer, "YES")
           x := 1
           half := n / 2
           for i := 0; i < half; i++ {
               v1 := make([]int, k)
               v2 := make([]int, k)
               for j := 0; j < k; j++ {
                   v1[j] = x
                   x++
                   v2[j] = x
                   x++
               }
               for j := 0; j < k; j++ {
                   fmt.Fprint(writer, v1[j], " ")
               }
               writer.WriteByte('\n')
               for j := 0; j < k; j++ {
                   fmt.Fprint(writer, v2[j], " ")
               }
               writer.WriteByte('\n')
           }
       }
   }
}
