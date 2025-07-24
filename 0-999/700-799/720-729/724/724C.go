package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

// extendedGcd returns (g, x, y) such that a*x + b*y = g = gcd(a,b)
func extendedGcd(a, b int64) (int64, int64, int64) {
   if b == 0 {
       return a, 1, 0
   }
   g, x1, y1 := extendedGcd(b, a%b)
   x := y1
   y := x1 - (a/b)*y1
   return g, x, y
}

// modInv returns modular inverse of a mod m, assuming gcd(a,m)==1
func modInv(a, m int64) int64 {
   _, x, _ := extendedGcd(a, m)
   x %= m
   if x < 0 {
       x += m
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, k int64
   fmt.Fscan(reader, &n, &m, &k)
   pts := make([][2]int64, k)
   for i := range pts {
       fmt.Fscan(reader, &pts[i][0], &pts[i][1])
   }
   // periods
   N := 2 * n
   M := 2 * m
   // gcd and lcm base
   g := gcd(N, M)
   // p = N/g, q = M/g
   p := N / g
   q := M / g
   invP := modInv(p%q, q)
   // process each sensor
   const INF = int64(9e18)
   res := make([]int64, k)
   for i, pt := range pts {
       x, y := pt[0], pt[1]
       best := INF
       // possible reflections
       vx := [2]int64{x, N - x}
       vy := [2]int64{y, M - y}
       for _, dx := range vx {
           for _, dy := range vy {
               // solve t ≡ dx mod N, t ≡ dy mod M
               diff := dy - dx
               if diff%g != 0 {
                   continue
               }
               // k * N + dx = t ≡ dy (mod M) -> k*(N/g) ≡ (dy-dx)/g mod (M/g)
               d := diff / g
               // ensure mod positive
               dMod := d % q
               if dMod < 0 {
                   dMod += q
               }
               k0 := (dMod * invP) % q
               t := dx + N*k0
               if t < best {
                   best = t
               }
           }
       }
       if best == INF {
           res[i] = -1
       } else {
           res[i] = best
       }
   }
   // output
   for i, v := range res {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}
