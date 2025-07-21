package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

// isqrt returns floor of square root of n
func isqrt(n int64) int64 {
   if n <= 0 {
       return 0
   }
   x := int64(math.Sqrt(float64(n)))
   for (x+1)*(x+1) <= n {
       x++
   }
   for x*x > n {
       x--
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var k int64
   _, _ = fmt.Fscan(reader, &k)
   // R is max distance for cell centers
   R := k - 1
   // T = 4*R^2 for transformed inequality X^2 + 3Y^2 <= T
   T := 4 * R * R
   // compute maximum Y for which 3*Y^2 <= T
   Ymax := isqrt(T / 3)
   var ans int64
   for y := int64(0); y <= Ymax; y++ {
       t := T - 3*y*y
       if t < 0 {
           continue
       }
       xmax := isqrt(t)
       var cnt int64
       if y%2 == 0 {
           // X even
           cnt = 2*(xmax/2) + 1
       } else {
           // X odd
           cnt = 2*((xmax+1)/2)
       }
       if y == 0 {
           ans += cnt
       } else {
           ans += 2 * cnt
       }
   }
   fmt.Println(ans)
}
