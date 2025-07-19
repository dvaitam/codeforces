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
   fmt.Fscan(reader, &n)

   r := 0
   for (1 << r) < n {
       r++
   }
   fmt.Fprintln(writer, r)

   for i := 0; i < r; i++ {
       var indices []int
       for j := 0; j < n; j++ {
           if (j>>i)&1 == 1 {
               indices = append(indices, j+1)
           }
       }
       fmt.Fprint(writer, len(indices))
       for _, v := range indices {
           fmt.Fprint(writer, " ", v)
       }
       fmt.Fprintln(writer)
   }
}
