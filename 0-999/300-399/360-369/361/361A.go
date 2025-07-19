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

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   put := k - n + 1

   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if i == j {
               fmt.Fprint(writer, put)
           } else {
               fmt.Fprint(writer, 1)
           }
           if j < n-1 {
               fmt.Fprint(writer, " ")
           }
       }
       fmt.Fprint(writer, '\n')
   }
}
