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

   var sumA, sumB, sumC int64
   for i := 0; i < n; i++ {
       var x int64
       fmt.Fscan(reader, &x)
       sumA += x
   }
   for i := 0; i < n-1; i++ {
       var x int64
       fmt.Fscan(reader, &x)
       sumB += x
   }
   for i := 0; i < n-2; i++ {
       var x int64
       fmt.Fscan(reader, &x)
       sumC += x
   }

   missing1 := sumA - sumB
   missing2 := sumB - sumC
   fmt.Fprintln(writer, missing1, missing2)
}
