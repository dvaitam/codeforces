package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   const maxZ = 110
   a := make([]int, n+1)
   p := make([][]float64, n+1)
   for i := 1; i <= n; i++ {
       p[i] = make([]float64, maxZ)
   }
   var ans float64
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] < maxZ {
           p[i][a[i]] = 1.0
       }
       if a[i] == 0 {
           ans += 1.0
       }
   }
   var q int
   fmt.Fscan(reader, &q)
   tmp := make([]float64, maxZ)
   for qi := 0; qi < q; qi++ {
       var u, v, k int
       fmt.Fscan(reader, &u, &v, &k)
       for j := 0; j < k; j++ {
           ans -= p[u][0]
           au := a[u]
           denom := float64(au)
           for z := 0; z <= au && z+1 < maxZ; z++ {
               tmp[z] = (p[u][z]*(denom-float64(z)) + p[u][z+1]*float64(z+1)) / denom
           }
           for z := 0; z <= au && z < maxZ; z++ {
               p[u][z] = tmp[z]
           }
           a[u]--
           ans += p[u][0]
       }
       a[v] += k
       fmt.Fprintf(writer, "%.10f\n", ans)
   }
}
