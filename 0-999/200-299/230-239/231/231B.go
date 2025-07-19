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

   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   d := make([]int, n)
   for i := 0; i < n; i++ {
       if m > 0 {
           d[i] = k
       } else {
           d[i] = 1
       }
       m = d[i] - m
   }
   // adjust last element by remaining m
   d[n-1] -= m
   if d[n-1] > 0 && d[n-1] <= k {
       for i := 0; i < n; i++ {
           if i > 0 {
               fmt.Fprint(writer, " ")
           }
           fmt.Fprint(writer, d[i])
       }
   } else {
       fmt.Fprint(writer, -1)
   }
}
