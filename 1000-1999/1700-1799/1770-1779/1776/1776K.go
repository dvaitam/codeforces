package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

const (
   B    = 100000
   maxK = 100
)

var (
   s [][]float64
   a []float64
)

// po computes n^{-k}
func po(n int64, k int) float64 {
   res := 1.0
   nn := float64(n)
   for i := 0; i < k; i++ {
       res /= nn
   }
   return res
}

// ph computes integral approximation for sum j^{-k}
func ph(n int64, k int) float64 {
   if n == 0 {
       return 0
   }
   nn := float64(n)
   switch k {
   case 0:
       return nn
   case 1:
       return math.Log(nn) + 0.5/nn
   default:
       return -1.0/float64(k-1) * po(n, k-1)
   }
}

// ca computes prefix sum of j^{-k} up to n
func ca(n int64, k int) float64 {
   if n <= B {
       return s[k][n]
   }
   return s[k][B] + (ph(n, k) - ph(B, k))
}

// cal computes probability distribution c for difference n
func cal(n int64, c []float64) {
   // adjust n: probability distribution uses n = n-1
   n--
   var f [maxK + 1]float64
   for i := 1; i <= maxK; i++ {
       val := ca(n, i)
       if i%2 == 0 {
           val = -val
       }
       f[i] = val / float64(i)
   }
   c[0] = 1.0
   for i := 1; i <= maxK; i++ {
       var res float64
       for j := 0; j < i; j++ {
           res += c[j] * float64(i-j) * f[i-j]
       }
       c[i] = res / float64(i)
   }
   // shift and normalize by n+1
   for i := maxK; i >= 1; i-- {
       c[i] = c[i-1] / float64(n+1)
   }
   // zero out invalid entries
   limit := int(n + 2)
   if limit > maxK {
       limit = maxK
   }
   for i := limit + 1; i <= maxK; i++ {
       c[i] = 0
   }
   c[0] = 0
}

func main() {
   // precompute a and s
   a = make([]float64, B+1)
   s = make([][]float64, maxK+1)
   for i := 0; i <= maxK; i++ {
       s[i] = make([]float64, B+1)
   }
   for j := 1; j <= B; j++ {
       a[j] = 1.0
       s[0][j] = float64(j)
   }
   for i := 1; i <= maxK; i++ {
       for j := 1; j <= B; j++ {
           a[j] /= float64(j)
           s[i][j] = s[i][j-1] + a[j]
       }
   }

   in := bufio.NewReader(os.Stdin)
   var T int64
   var m int
   fmt.Fscan(in, &T, &m)
   b := make([]int64, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &b[i])
   }
   if m == 1 {
       fmt.Printf("1\n")
       return
   }
   // compute c for each researcher
   c := make([][]float64, m)
   for i := 0; i < m; i++ {
       c[i] = make([]float64, maxK+1)
       if b[i] == T {
           // zero difference: trivial
           continue
       }
       cal(T-b[i], c[i])
   }
   // compute answers
   ans := make([]float64, m)
   for i := 0; i < m; i++ {
       nw := make([]float64, m)
       for j := 0; j < m; j++ {
           nw[j] = 1.0
       }
       var resi float64
       for year := 1; year <= maxK; year++ {
           // compute probability researcher i wins at this year
           prod := 1.0
           for k := 0; k < m; k++ {
               if k == i {
                   continue
               }
               prod *= nw[k]
           }
           prob := c[i][year] * prod
           resi += prob
           // update survival probabilities
           for k := 0; k < m; k++ {
               nw[k] -= c[k][year]
               if nw[k] < 0 {
                   nw[k] = 0
               }
           }
       }
       ans[i] = resi
   }
   // output
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i, v := range ans {
       if i > 0 {
           fmt.Fprint(out, ' ')
       }
       fmt.Fprintf(out, "%.12f", v)
   }
  fmt.Fprint(out, '\n')
}
