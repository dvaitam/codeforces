package main

import "fmt"

func main() {
   var n, k, m int
   if _, err := fmt.Scan(&n, &k, &m); err != nil {
       return
   }
   // precompute powers of 10 modulo k
   pow10 := make([]int, n)
   pow10[0] = 1 % k
   for i := 1; i < n; i++ {
       pow10[i] = pow10[i-1] * 10 % k
   }
   // dp arrays: dpCurr[r][f] = ways with current prefix mod r and flag f (found valid suffix)
   dpCurr := make([][2]int, k)
   dpNext := make([][2]int, k)
   dpCurr[0][0] = 1 % m
   for i := 0; i < n; i++ {
       // reset next
       for r := 0; r < k; r++ {
           dpNext[r][0] = 0
           dpNext[r][1] = 0
       }
       // last digit (most significant of original) must be non-zero
       start := 0
       if i+1 == n {
           start = 1
       }
       for r := 0; r < k; r++ {
           for flag := 0; flag < 2; flag++ {
               v := dpCurr[r][flag]
               if v == 0 {
                   continue
               }
               for d := start; d <= 9; d++ {
                   newR := (d*pow10[i] + r) % k
                   newFlag := flag
                   // if this prefix (suffix of original) is divisible and has no leading zero
                   if flag == 0 && newR == 0 && d != 0 {
                       newFlag = 1
                   }
                   dpNext[newR][newFlag] += v
                   if dpNext[newR][newFlag] >= m {
                       dpNext[newR][newFlag] %= m
                   }
               }
           }
       }
       dpCurr, dpNext = dpNext, dpCurr
   }
   // sum all states where we found at least one valid suffix
   ans := 0
   for r := 0; r < k; r++ {
       ans += dpCurr[r][1]
       if ans >= m {
           ans %= m
       }
   }
   fmt.Println(ans)
}
