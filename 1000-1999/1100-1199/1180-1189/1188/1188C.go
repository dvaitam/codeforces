package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const mod = 998244353

func add(a, b int) int {
   a += b
   if a >= mod {
       a -= mod
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)
   // maximum possible diff
   maxDiff := a[n-1] - a[0]
   // dp[c][i]: number of ways to choose c elements ending at position i (1-based)
   dp := make([][]int, k+1)
   for i := range dp {
       dp[i] = make([]int, n)
   }
   L := make([]int, n)
   pref := make([]int, n+1)
   total := 0
   // threshold t from 1 to maxDiff
   for t := 1; t <= maxDiff; t++ {
       // compute L[i]: largest j < i such that a[i] - a[j] >= t
       j := 0
       for i := 0; i < n; i++ {
           for j < i && a[i]-a[j] >= t {
               j++
           }
           L[i] = j - 1
           if L[i] < 0 {
               L[i] = -1
           }
       }
       // dp for this t
       // c = 1
       for i := 0; i < n; i++ {
           dp[1][i] = 1
       }
       // c from 2 to k
       for c := 2; c <= k; c++ {
           // prefix sums of dp[c-1]
           pref[0] = 0
           for i := 0; i < n; i++ {
               pref[i+1] = add(pref[i], dp[c-1][i])
           }
           for i := 0; i < n; i++ {
               li := L[i]
               if li >= 0 {
                   dp[c][i] = pref[li+1]
               } else {
                   dp[c][i] = 0
               }
           }
       }
       // sum dp[k]
       sumk := 0
       for i := 0; i < n; i++ {
           sumk = add(sumk, dp[k][i])
       }
       total = add(total, sumk)
   }
   fmt.Fprintln(writer, total)
}
