package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   p := make([]float64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &p[i])
   }
   sort.Float64s(p)
   const eps = 1e-9
   if n > 0 && math.Abs(p[n-1]-1.0) < eps {
       fmt.Printf("%.15f\n", 1.0)
       return
   }
   // prefix and suffix products of (1-p)
   prod := make([]float64, n+1)
   revprod := make([]float64, n+1)
   prod[0] = 1.0
   for i := 0; i < n; i++ {
       prod[i+1] = prod[i] * (1.0 - p[i])
   }
   revprod[n] = 1.0
   for i := n - 1; i >= 0; i-- {
       revprod[i] = revprod[i+1] * (1.0 - p[i])
   }
   best := 0.0
   // choose i smallest and (n-j) largest to remove
   for i := 0; i <= n; i++ {
       for j := n; j >= i; j-- {
           pr := prod[i] * revprod[j]
           cnd := 0.0
           for k := 0; k < i; k++ {
               if 1.0-p[k] > eps {
                   cnd += p[k] * pr / (1.0 - p[k])
               }
           }
           for k := n - 1; k >= j; k-- {
               if 1.0-p[k] > eps {
                   cnd += p[k] * pr / (1.0 - p[k])
               }
           }
           if cnd > best {
               best = cnd
           }
       }
   }
   fmt.Printf("%.15f\n", best)
}
