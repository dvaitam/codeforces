package main

import (
   "fmt"
   "math"
)

const maxn = 30

var (
   nf, ne, ns int
   df, de   int
   rf, re_, rs float64
   U         [maxn]bool
   a         [maxn]int
   b         [maxn]float64
   ans       float64
)

// Len returns overlap length between [a,b] and [c,d]
func Len(a1, b1, c1, d1 float64) float64 {
   lo := math.Max(a1, c1)
   hi := math.Min(b1, d1)
   if hi > lo {
       return hi - lo
   }
   return 0.0
}

// calc computes total damage for current freezing placement U and positions a[0:ns]
func calc() float64 {
   // base damage without electric bonuses
   Fc := 2*float64(nf)*rf*float64(df) + 2*float64(ne)*re_*float64(de)
   m := 0
   for i := 0; i < nf+ne+ns; i++ {
       if !U[i] {
           var Df, DeF float64
           xi := float64(i) / 2.0
           for j := 0; j < ns; j++ {
               // slowed damage segments due to freezing
               Df += float64(df) * Len(xi-rf, xi+rf, float64(a[j])-rs, float64(a[j])+rs)
               DeF += float64(de) * Len(xi-re_, xi+re_, float64(a[j])-rs, float64(a[j])+rs)
           }
           Fc += Df
           b[m] = DeF - Df
           m++
       }
   }
   // choose top ne electric bonuses: simple descending sort
   for i := 0; i < m; i++ {
       for j := i + 1; j < m; j++ {
           if b[j] > b[i] {
               b[i], b[j] = b[j], b[i]
           }
       }
   }
   for i := 0; i < ne && i < m; i++ {
       Fc += b[i]
   }
   return Fc
}

// dfs tries all assignments of freezing towers
func dfs(x, y int) {
   if nf+ne+y < x {
       return
   }
   if x == nf+ne+ns {
       val := calc()
       if val > ans {
           ans = val
       }
       return
   }
   // no freezing at position x
   U[x] = false
   dfs(x+1, y)
   // place freezing if possible
   if y < ns && (x%2 == 0 || U[x-1]) {
       U[x] = true
       a[y] = x / 2
       dfs(x+1, y+1)
   }
}

func main() {
   var rfi, rei, rsi float64
   fmt.Scan(&nf, &ne, &ns)
   fmt.Scan(&rfi, &rei, &rsi)
   fmt.Scan(&df, &de)
   // adjust radii for y-offset of 1
   rf = math.Sqrt(rfi*rfi - 1)
   re_ = math.Sqrt(rei*rei - 1)
   rs = math.Sqrt(rsi*rsi - 1)
   dfs(0, 0)
   fmt.Printf("%.10f\n", ans)
}
