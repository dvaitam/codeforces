package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func isqrt(x int64) int64 {
   r := int64(math.Sqrt(float64(x)))
   for (r+1)*(r+1) <= x {
      r++
   }
   for r*r > x {
      r--
   }
   return r
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
      t--
      var n int64
      fmt.Fscan(reader, &n)
      x := (1 + isqrt(1+8*n)) / 2
      for x*(x-1)/2 > n {
         x--
      }
      tri := x * (x - 1) / 2
      res := x + (n - tri)
      fmt.Fprintln(writer, res)
   }
}
