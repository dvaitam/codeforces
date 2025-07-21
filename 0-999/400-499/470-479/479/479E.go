package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, a, b, k int
   fmt.Fscan(in, &n, &a, &b, &k)
   // dp[x]: ways to be at floor x
   dp := make([]int, n+2)
   dp[a] = 1
   for step := 1; step <= k; step++ {
       // prefix sums
       pref := make([]int, n+2)
       for i := 1; i <= n; i++ {
           pref[i] = pref[i-1] + dp[i]
           if pref[i] >= mod {
               pref[i] -= mod
           }
       }
       newdp := make([]int, n+2)
       for x := 1; x <= n; x++ {
           if x == b {
               continue
           }
           d := b - x
           if d < 0 {
               d = -d
           }
           // allowed y: |y-x| < d => y in [x-d+1, x+d-1]
           l := x - d + 1
           if l < 1 {
               l = 1
           }
           r := x + d - 1
           if r > n {
               r = n
           }
           total := pref[r] - pref[l-1]
           if total < 0 {
               total += mod
           }
           // exclude staying at x
           total = total - dp[x]
           if total < 0 {
               total += mod
           }
           newdp[x] = total
       }
       dp = newdp
   }
   ans := 0
   for i := 1; i <= n; i++ {
       ans += dp[i]
       if ans >= mod {
           ans -= mod
       }
   }
   fmt.Fprintln(out, ans)
}
