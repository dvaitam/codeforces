package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int64, n+1)
   b := make([]int64, n+1)
   c := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &b[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &c[i])
   }
   const INF = int64(1e18)
   dp0 := [2]int64{-INF, -INF}
   dp0[0] = 0
   // DP over edges x_i: x_i = 1 if i fed before i+1
   for i := 1; i < n; i++ {
       dp1 := [2]int64{-INF, -INF}
       for px := 0; px <= 1; px++ {
           base := dp0[px]
           if base <= -INF/2 {
               continue
           }
           for xi := 0; xi <= 1; xi++ {
               // xi is x_i
               count := px + (1 - xi)
               var joy int64
               switch count {
               case 0:
                   joy = a[i]
               case 1:
                   joy = b[i]
               default:
                   joy = c[i]
               }
               val := base + joy
               if val > dp1[xi] {
                   dp1[xi] = val
               }
           }
       }
       dp0 = dp1
   }
   // last hare n has only left neighbor
   // if x_{n-1} = 0 -> neighbor later -> 0 full -> a[n]
   // if x_{n-1} = 1 -> neighbor earlier -> 1 full -> b[n]
   ans := dp0[0] + a[n]
   if v := dp0[1] + b[n]; v > ans {
       ans = v
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprintln(out, ans)
}
