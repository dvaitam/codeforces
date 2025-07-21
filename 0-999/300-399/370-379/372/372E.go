package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   if a < 0 {
       return -a
   }
   return a
}

func modPow(a, e int64) int64 {
   res := int64(1)
   a %= mod
   for e > 0 {
       if e&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       e >>= 1
   }
   return res
}

type pair struct{
   x, y int64
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   xs := make([]int64, n)
   ys := make([]int64, n)
   for i := 0; i < n; i++ {
       var ai, bi, ci, di int64
       fmt.Fscan(in, &ai, &bi, &ci, &di)
       // Xi = ai/bi, Yi = ci/di; scale by bi*di
       xs[i] = ai * di
       ys[i] = ci * bi
   }
   counts := make(map[pair]int)
   for i := 0; i < n; i++ {
       xi, yi := xs[i], ys[i]
       for j := i + 1; j < n; j++ {
           xj, yj := xs[j], ys[j]
           // compute a*b = (xi+iyi)*(xj+iyj) = dot + i*cross
           dot := xi*xj - yi*yj
           cross := xi*yj + yi*xj
           // normalize
           g := gcd(dot, cross)
           if g != 0 {
               dot /= g
               cross /= g
           }
           key := pair{dot, cross}
           counts[key]++
       }
   }
   var ans int64 = 0
   for _, m := range counts {
       if m >= 2 {
           // sum of C(m, k) for k>=2 = 2^m - 1 - m
           add := (modPow(2, int64(m)) - 1 - int64(m)) % mod
           if add < 0 {
               add += mod
           }
           ans = (ans + add) % mod
       }
   }
   fmt.Println(ans)
}
