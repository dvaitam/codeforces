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
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       p := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &p[i])
       }
       // output reversed permutation
       for i := n - 1; i >= 0; i-- {
           fmt.Fprint(writer, p[i])
           if i > 0 {
               fmt.Fprint(writer, ' ')
           }
       }
       fmt.Fprint(writer, '\n')
   }
}
