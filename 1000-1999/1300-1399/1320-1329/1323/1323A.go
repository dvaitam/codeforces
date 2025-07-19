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
       ev, od := 0, 0
       var eva, oda []int
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
           if a[i]%2 == 0 {
               ev++
               eva = append(eva, i+1)
           } else {
               od++
               oda = append(oda, i+1)
           }
       }
       if n == 1 && od == 1 {
           fmt.Fprintln(writer, -1)
       } else if ev >= 1 {
           fmt.Fprintln(writer, 1)
           fmt.Fprintln(writer, eva[0])
       } else {
           fmt.Fprintln(writer, 2)
           fmt.Fprintf(writer, "%d %d\n", oda[0], oda[1])
       }
   }
}
