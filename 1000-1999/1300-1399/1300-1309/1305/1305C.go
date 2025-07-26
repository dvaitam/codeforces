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

   var n, m int
   fmt.Fscan(reader, &n, &m)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // if more numbers than modulo, pigeonhole implies zero
   if n > m {
       fmt.Fprint(writer, 0)
       return
   }
   result := 1
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           diff := a[i] - a[j]
           if diff < 0 {
               diff = -diff
           }
           result = result * (diff % m) % m
       }
   }
   fmt.Fprint(writer, result)
}
