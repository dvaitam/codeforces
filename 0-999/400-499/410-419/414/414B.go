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

   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }

   // dp for sequences: use rolling arrays prev and cur
   prev := make([]int, n+1)
   for i := 1; i <= n; i++ {
       prev[i] = 1
   }

   cur := make([]int, n+1)
   for length := 2; length <= k; length++ {
       // reset cur
       for i := 1; i <= n; i++ {
           cur[i] = 0
       }
       // for each possible previous value x, propagate to multiples
       for x := 1; x <= n; x++ {
           v := prev[x]
           if v == 0 {
               continue
           }
           for m := x; m <= n; m += x {
               cur[m] += v
               if cur[m] >= mod {
                   cur[m] -= mod
               }
           }
       }
       // swap prev and cur
       prev, cur = cur, prev
   }

   // sum up sequences of length k
   ans := 0
   for i := 1; i <= n; i++ {
       ans += prev[i]
       if ans >= mod {
           ans -= mod
       }
   }
   fmt.Fprintln(out, ans)
}
