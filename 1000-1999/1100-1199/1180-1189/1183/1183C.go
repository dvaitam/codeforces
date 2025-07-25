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

   var q int
   if _, err := fmt.Fscan(reader, &q); err != nil {
       return
   }
   for i := 0; i < q; i++ {
       var k, n, a, b int64
       fmt.Fscan(reader, &k, &n, &a, &b)
       // If even all turns with light use drain too much
       if k <= n*b {
           fmt.Fprintln(writer, -1)
           continue
       }
       // Maximum heavy plays x such that k - n*b - x*(a-b) > 0
       // x <= (k - n*b - 1) / (a - b)
       maxX := (k - n*b - 1) / (a - b)
       if maxX > n {
           maxX = n
       }
       if maxX < 0 {
           maxX = 0
       }
       fmt.Fprintln(writer, maxX)
   }
}
