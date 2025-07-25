package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int64
   fmt.Fscan(reader, &n, &k)

   // Solve for number of put operations m using quadratic: m^2 +3m -2(n+k) = 0
   // Discriminant t = 9 + 8(n+k)
   t := 8*(n + k) + 9
   s := int64(math.Sqrt(float64(t)))
   if s*s != t {
       // Adjust in case of rounding
       if (s+1)*(s+1) == t {
           s++
       }
   }
   // m = (s - 3) / 2
   m := (s - 3) / 2
   // Number of eaten candies x = total moves - put moves
   x := n - m
   fmt.Fprintln(writer, x)
}
