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
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       res := make([]int, 0, 2*n)
       res = append(res, a[0])
       for i := 1; i < n; i++ {
           if a[i] < a[i-1] {
               res = append(res, a[i])
           }
           res = append(res, a[i])
       }
       fmt.Fprintln(writer, len(res))
       for i, v := range res {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       writer.WriteByte('\n')
   }
}
