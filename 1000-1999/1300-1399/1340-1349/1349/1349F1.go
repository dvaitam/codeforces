package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // Precompute dpR: dpR[l][y] = number of suffix sequences of length l starting with max=y
   // l from 0..n-1, y from 1..n+1
   dpR := make([][]int, n)
   // dpR[0]
   dpR[0] = make([]int, n+2)
   for y := 1; y <= n+1; y++ {
       dpR[0][y] = 1
   }
   for l := 1; l < n; l++ {
       prev := dpR[l-1]
       cur := make([]int, n+2)
       // dpR[l][y] = y*prev[y] + prev[y+1]
       for y := n; y >= 1; y-- {
           v := (int64(y)*int64(prev[y]) + int64(prev[y+1])) % mod
           cur[y] = int(v)
       }
       dpR[l] = cur
   }
   // S_prev holds S(i-1,m) for current i
   S_prev := make([]int, n+2)
   S_prev[0] = 1
   S_cur := make([]int, n+2)
   // sum1 and sum2 accumulators
   sum1 := make([]int, n+2)
   sum2 := make([]int, n+2)
   // process for each i
   for i := 1; i <= n; i++ {
       // dpR row for suffix length
       l := n - i
       row := dpR[l]
       // suffix sum and accumulate
       var suffix int
       for m := n; m >= 1; m-- {
           // V = S_prev[m] * row[m]
           v := int64(S_prev[m]) * int64(row[m]) % mod
           suffix = int((int64(suffix) + v) % mod)
           // sum1[m] += suffix
           sum1[m] = int((int64(sum1[m]) + int64(suffix)) % mod)
           // sum2[m] += S_prev[m-1] * row[m]
           v2 := int64(S_prev[m-1]) * int64(row[m]) % mod
           sum2[m] = int((int64(sum2[m]) + v2) % mod)
       }
       // update S_cur = S(i,m)
       S_cur[0] = 0
       for m := 1; m <= i; m++ {
           S_cur[m] = int((int64(S_prev[m])*int64(m) + int64(S_prev[m-1])) % mod)
       }
       // clear rest
       for m := i + 1; m <= n; m++ {
           S_cur[m] = 0
       }
       // swap
       S_prev, S_cur = S_cur, S_prev
   }
   // output answers
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for k := 1; k <= n; k++ {
       ans := sum1[k] + sum2[k]
       if ans >= mod {
           ans -= mod
       }
       fmt.Fprintf(out, "%d", ans)
       if k < n {
           out.WriteByte(' ')
       }
   }
   out.WriteByte('\n')
}
