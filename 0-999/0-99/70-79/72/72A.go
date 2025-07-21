package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // generate primes up to n
   isPrime := make([]bool, n+1)
   for i := 2; i <= n; i++ {
       isPrime[i] = true
   }
   for i := 2; i*i <= n; i++ {
       if isPrime[i] {
           for j := i * i; j <= n; j += i {
               isPrime[j] = false
           }
       }
   }
   // build candidates: primes in descending order, then 1
   primes := make([]int, 0, n/10)
   for i := 2; i <= n; i++ {
       if isPrime[i] {
           primes = append(primes, i)
       }
   }
   // vals: primes descending
   mpr := len(primes)
   vals := make([]int, 0, mpr+1)
   for i := mpr - 1; i >= 0; i-- {
       vals = append(vals, primes[i])
   }
   // include 1 as smallest
   if n >= 1 {
       vals = append(vals, 1)
   }
   m := len(vals)
   // dp[i][s]: using vals[i:], can we sum to s
   dp := make([][]bool, m+1)
   for i := range dp {
       dp[i] = make([]bool, n+1)
   }
   dp[m][0] = true
   for i := m - 1; i >= 0; i-- {
       vi := vals[i]
       dpi1 := dp[i+1]
       dpi := dp[i]
       for s := 0; s <= n; s++ {
           if dpi1[s] {
               dpi[s] = true
           } else if s >= vi && dpi1[s-vi] {
               dpi[s] = true
           }
       }
   }
   if !dp[0][n] {
       fmt.Println(0)
       return
   }
   // greedy pick lexicographically largest
   res := make([]int, 0, m)
   sum := n
   idx := 0
   for sum > 0 {
       picked := false
       for i := idx; i < m; i++ {
           v := vals[i]
           if v > sum {
               continue
           }
           if dp[i+1][sum-v] {
               res = append(res, v)
               sum -= v
               idx = i + 1
               picked = true
               break
           }
       }
       if !picked {
           // no solution (should not happen)
           fmt.Println(0)
           return
       }
   }
   // output result
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i, v := range res {
       if i > 0 {
           w.WriteByte(' ')
       }
       fmt.Fprint(w, v)
   }
   w.WriteByte('\n')
}
