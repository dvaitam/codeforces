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

   var n, v int
   if _, err := fmt.Fscan(reader, &n, &v); err != nil {
       return
   }
   var res []int
   for i := 1; i <= n; i++ {
       var k int
       fmt.Fscan(reader, &k)
       ok := false
       for j := 0; j < k; j++ {
           var s int
           fmt.Fscan(reader, &s)
           if s < v {
               ok = true
           }
       }
       if ok {
           res = append(res, i)
       }
   }
   // Output
   fmt.Fprintln(writer, len(res))
   for idx, seller := range res {
       if idx > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, seller)
   }
   // End the second line (empty if no sellers)
   fmt.Fprintln(writer)
}
