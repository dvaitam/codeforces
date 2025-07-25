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

   var n, m, _h int
   if _, err := fmt.Fscan(reader, &n, &m, &_h); err != nil {
       return
   }
   a := make([]int, m)
   for j := 0; j < m; j++ {
       fmt.Fscan(reader, &a[j])
   }
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   t := make([][]int, n)
   for i := 0; i < n; i++ {
       t[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &t[i][j])
       }
   }
   // construct heights
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           var hval int
           if t[i][j] != 0 {
               // assign min of front and left view
               if a[j] < b[i] {
                   hval = a[j]
               } else {
                   hval = b[i]
               }
           }
           if j > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, hval)
       }
       writer.WriteByte('\n')
   }
}
