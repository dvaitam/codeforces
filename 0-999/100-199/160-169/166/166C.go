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

   var n, x int
   fmt.Fscan(reader, &n, &x)
   arr := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &arr[i])
   }

   L, E := 0, 0
   for _, v := range arr {
       if v < x {
           L++
       } else if v == x {
           E++
       }
   }

   k := 0
   for {
       N := n + k
       pos := (N + 1) / 2
       if pos > L && pos <= L+E+k {
           fmt.Fprintln(writer, k)
           break
       }
       k++
   }
}
