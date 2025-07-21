package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   // compute deltas and sum of positive and negative deltas
   deltas := make([]int, n)
   negSum, posSum := 0, 0
   for i := 0; i < n; i++ {
       del := a[i] - k*b[i]
       deltas[i] = del
       if del < 0 {
           negSum += del
       } else {
           posSum += del
       }
   }
   // dp index offset
   base := -negSum
   size := posSum + base + 1
   // dp[j] = max sum(a) for sum(delta) == j-base, -1 unreachable
   dp := make([]int, size)
   for i := 0; i < size; i++ {
       dp[i] = -1
   }
   dp[base] = 0
   // 0-1 knapsack over deltas
   for i := 0; i < n; i++ {
       del := deltas[i]
       ai := a[i]
       if del >= 0 {
           // iterate from high to low
           for j := size - 1 - del; j >= 0; j-- {
               if dp[j] >= 0 {
                   nj := j + del
                   v := dp[j] + ai
                   if v > dp[nj] {
                       dp[nj] = v
                   }
               }
           }
       } else {
           // del < 0, iterate from low to high
           for j := -del; j < size; j++ {
               if dp[j] >= 0 {
                   nj := j + del
                   v := dp[j] + ai
                   if v > dp[nj] {
                       dp[nj] = v
                   }
               }
           }
       }
   }
   // result at sum(delta)==0 => index base
   res := dp[base]
   if res > 0 {
       fmt.Println(res)
   } else {
       fmt.Println(-1)
   }
}
