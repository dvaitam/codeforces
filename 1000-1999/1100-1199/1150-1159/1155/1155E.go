package main

import (
   "bufio"
   "fmt"
   "os"
)

const M int64 = 1000003

func modpow(x, e int64) int64 {
   res := int64(1)
   for e > 0 {
       if e&1 == 1 {
           res = res * x % M
       }
       x = x * x % M
       e >>= 1
   }
   return res
}

func modinv(a int64) int64 {
   return modpow((a%M+M)%M, M-2)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // Query points x = 0..10
   n := 10
   xs := make([]int64, n+1)
   ys := make([]int64, n+1)
   for i := 0; i <= n; i++ {
       xs[i] = int64(i)
       fmt.Fprintf(writer, "? %d\n", i)
       writer.Flush()
       var y int64
       if _, err := fmt.Fscan(reader, &y); err != nil {
           os.Exit(0)
       }
       ys[i] = (y%M + M) % M
   }

   // Compute polynomial coefficients via Lagrange interpolation
   // f(x) = sum ys[i] * Li(x)
   // Li(x) = Prod_{j!=i}(x - xs[j]) / Prod_{j!=i}(xs[i] - xs[j])
   // Build coefficients a[0..n]
   a := make([]int64, n+1)
   // Precompute denominators
   den := make([]int64, n+1)
   for i := 0; i <= n; i++ {
       d := int64(1)
       for j := 0; j <= n; j++ {
           if j == i {
               continue
           }
           d = d * ((xs[i] - xs[j] + M) % M) % M
       }
       den[i] = modinv(d)
   }
   // For each basis polynomial
   for i := 0; i <= n; i++ {
       // build P(x) = Prod_{j!=i} (x - xs[j])
       // poly init to 1
       poly := make([]int64, 1)
       poly[0] = 1
       for j := 0; j <= n; j++ {
           if j == i {
               continue
           }
           // multiply poly by (x - xs[j])
           nj := len(poly)
           newp := make([]int64, nj+1)
           for k := 0; k < nj; k++ {
               // coeff for x^{k+1}
               newp[k+1] = (newp[k+1] + poly[k]) % M
               // coeff for x^{k}
               newp[k] = (newp[k] - poly[k]*xs[j]) % M
           }
           poly = newp
       }
       // scale by ys[i] * den[i]
       factor := ys[i] * den[i] % M
       for k := 0; k < len(poly); k++ {
           a[k] = (a[k] + poly[k]*factor) % M
       }
   }
   // Normalize coefficients to [0,M)
   for k := 0; k <= n; k++ {
       a[k] = (a[k] + M) % M
   }

   // Search for root by brute force
   for x := int64(0); x < M; x++ {
       // evaluate f(x)
       res := int64(0)
       for k := n; k >= 0; k-- {
           res = (res*x + a[k]) % M
       }
       if res == 0 {
           fmt.Fprintf(writer, "! %d\n", x)
           writer.Flush()
           return
       }
   }
   // no root
   fmt.Fprint(writer, "! -1\n")
}
