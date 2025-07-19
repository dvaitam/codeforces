package main

import (
   "bufio"
   "fmt"
   "os"
)

// exgcd returns gcd(a,b) and x,y such that a*x + b*y = gcd
func exgcd(a, b int64) (g, x, y int64) {
   if b == 0 {
       return a, 1, 0
   }
   g, x1, y1 := exgcd(b, a%b)
   x = y1
   y = x1 - (a/b)*y1
   return
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m, dx, dy int64
   if _, err := fmt.Fscan(in, &n, &m, &dx, &dy); err != nil {
       return
   }
   // compute inverse of dx modulo n
   g, invDx, _ := exgcd(dx, n)
   if g != 1 {
       return
   }
   invDx = (invDx % n + n) % n
   // starting y-coordinate offset
   yCoord := (n - invDx) % n

   // count hits for each possible starting y
   a := make([]int, n)
   for i := int64(0); i < m; i++ {
       var xi, yi int64
       fmt.Fscan(in, &xi, &yi)
       // t = (xi * yCoord) % n * dy % n
       t := (xi*yCoord)%n * dy % n
       k := (yi + t) % n
       a[int(k)]++
   }
   // find maximum
   ans := 0
   for i := 1; i < int(n); i++ {
       if a[i] > a[ans] {
           ans = i
       }
   }
   fmt.Fprintf(out, "0 %d\n", ans)
}
