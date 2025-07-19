package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int, n+1)
   var t [3]int
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &t[0], &t[1], &t[2])
       var b [4]bool
       for j := 0; j < 3; j++ {
           b[a[t[j]]] = true
       }
       for j := 0; j < 3; j++ {
           if a[t[j]] == 0 {
               for k := 1; k <= 3; k++ {
                   if !b[k] {
                       a[t[j]] = k
                       b[k] = true
                       break
                   }
               }
           }
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, a[i])
   }
   writer.WriteByte('\n')
}
