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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // Maximum independent set on grid: ceil(n*n/2)
   total := n * n
   maxCoders := total / 2
   if total%2 != 0 {
       maxCoders++
   }
   fmt.Fprintln(writer, maxCoders)
   // Print grid: place on (i+j)%2==0
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if (i+j)%2 == 0 {
               writer.WriteByte('C')
           } else {
               writer.WriteByte('.')
           }
       }
       writer.WriteByte('\n')
   }
}
