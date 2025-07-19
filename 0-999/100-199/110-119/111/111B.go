package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   const maxV = 100000
   // smallest prime factor sieve
   spf := make([]int, maxV+1)
   for i := 2; i <= maxV; i++ {
       if spf[i] == 0 {
           spf[i] = i
           if i > maxV/i {
               continue
           }
           for j := i * i; j <= maxV; j += i {
               if spf[j] == 0 {
                   spf[j] = i
               }
           }
       }
   }
   // last occurrence of each divisor
   last := make([]int, maxV+1)
   for i := range last {
       last[i] = -1
   }
   for idx := 0; idx < n; idx++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       // factor x into primes and counts
       ps := make([]int, 0)
       cs := make([]int, 0)
       xx := x
       for xx > 1 {
           p := spf[xx]
           cnt := 0
           for xx%p == 0 {
               xx /= p
               cnt++
           }
           ps = append(ps, p)
           cs = append(cs, cnt)
       }
       l := idx - y
       r := idx - 1
       ans := 0
       // dfs to enumerate divisors
       var dfs func(step int, mul int)
       dfs = func(step int, mul int) {
           if step < 0 {
               if last[mul] < l {
                   ans++
               }
               last[mul] = r + 1
               return
           }
           cur := mul
           for i := 0; i <= cs[step]; i++ {
               dfs(step-1, cur)
               cur *= ps[step]
           }
       }
       dfs(len(ps)-1, 1)
       fmt.Fprintln(out, ans)
   }
}
