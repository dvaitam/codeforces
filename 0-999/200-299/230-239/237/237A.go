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
   counts := make([]int, 24*60)
   maxC := 0
   for i := 0; i < n; i++ {
       var h, m int
       fmt.Fscan(reader, &h, &m)
       idx := h*60 + m
       counts[idx]++
       if counts[idx] > maxC {
           maxC = counts[idx]
       }
   }
   fmt.Fprint(writer, maxC)
}
