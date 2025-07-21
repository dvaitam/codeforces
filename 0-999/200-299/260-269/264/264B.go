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

   // read input
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   arr := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &arr[i])
   }

   // max value for sieve
   const maxV = 100000
   spf := make([]int, maxV+1)
   for i := 2; i <= maxV; i++ {
       if spf[i] == 0 {
           for j := i; j <= maxV; j += i {
               if spf[j] == 0 {
                   spf[j] = i
               }
           }
       }
   }

   dp := make([]int, maxV+1)
   ans := 1
   for _, v := range arr {
       x := v
       // get unique prime factors
       primes := make([]int, 0, 8)
       for x > 1 {
           p := spf[x]
           primes = append(primes, p)
           for x%p == 0 {
               x /= p
           }
       }
       best := 1
       for _, p := range primes {
           if dp[p]+1 > best {
               best = dp[p] + 1
           }
       }
       for _, p := range primes {
           if dp[p] < best {
               dp[p] = best
           }
       }
       if best > ans {
           ans = best
       }
   }
   fmt.Fprintln(writer, ans)
}
