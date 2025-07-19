package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   A := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &A[i])
   }
   // dp[i][j][k]: minimal largest segment distance from city i to j with k refuels
   N := n
   NN := N * N
   size := NN * N
   dp := make([]int32, size)
   const inf32 = int32(1000000005)
   for i := 0; i < N; i++ {
       iNN := i * NN
       for j := i; j < N; j++ {
           ijBase := iNN + j*N
           // k = 0
           dp[ijBase] = int32(A[j] - A[i])
           s := i
           // k from 1 to j-i
           for k := 1; k <= j-i; k++ {
               // advance s to optimal
               for s < j-1 && dp[iNN + s*N + (k-1)] < int32(A[j]-A[s]) {
                   s++
               }
               // compute dp[i][j][k]
               var v int32 = inf32
               if s != i {
                   v1 := dp[iNN + (s-1)*N + (k-1)]
                   d1 := int32(A[j] - A[s-1])
                   if v1 > d1 {
                       v = v1
                   } else {
                       v = d1
                   }
               }
               v2 := dp[iNN + s*N + (k-1)]
               d2 := int32(A[j] - A[s])
               var t32 int32
               if v2 > d2 {
                   t32 = v2
               } else {
                   t32 = d2
               }
               if v > t32 {
                   v = t32
               }
               dp[ijBase + k] = v
           }
       }
   }
   var ans int64
   for q := 0; q < m; q++ {
       var s, t, c, r int
       fmt.Fscan(reader, &s, &t, &c, &r)
       s--
       t--
       maxRefuels := t - s
       if r > maxRefuels {
           r = maxRefuels
       }
       val := int64(dp[s*NN + t*N + r]) * int64(c)
       if val > ans {
           ans = val
       }
   }
   fmt.Println(ans)
}
