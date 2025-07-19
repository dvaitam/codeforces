package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   fac := make([]float64, m+1)
   fac[0] = 0.0
   for i := 1; i <= m; i++ {
       fac[i] = fac[i-1] + math.Log(float64(i))
   }
   C := func(a, b int) float64 {
       return fac[b] - fac[a] - fac[b-a]
   }
   var ans float64
   limit := math.Exp(1) // dummy to satisfy possible import, not used
   _ = limit
   for i := 0; i <= n; i++ {
       for j := 0; j <= n; j++ {
           t := i*n + j*n - i*j
           if t > k {
               break
           }
           tmp := C(k-t, m-t) + C(i, n) + C(j, n) - C(k, m)
           ans += math.Exp(tmp)
           if ans > 1e99 {
               fmt.Println("1e99")
               return
           }
       }
   }
   fmt.Printf("%.15f\n", ans)
}
