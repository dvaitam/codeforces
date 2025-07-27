package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var k int
   // read k
   if _, err := fmt.Fscan(in, &k); err != nil {
       return
   }
   // read fortunes F0..F5
   F := make([]int64, 6)
   for i := 0; i < 6; i++ {
       fmt.Fscan(in, &F[i])
   }
   // read q (always 1)
   var q int
   fmt.Fscan(in, &q)
   // process each query
   for qi := 0; qi < q; qi++ {
       var n int
       fmt.Fscan(in, &n)
       // for k>=5, leftover digits (<=2k>=10) always can fill any remaining weight, use knapsack on blocks
       if k >= 5 {
           cap := n
           // prepare bounded knapsack dp over block weights
           pow10 := [6]int{1, 10, 100, 1000, 10000, 100000}
           const inf = int64(-4e18)
           dp := make([]int64, cap+1)
           for i := 1; i <= cap; i++ {
               dp[i] = inf
           }
           dp[0] = 0
           // items: type j has weight w=3*10^j, profit=F[j], count=3*k
           for j := 0; j < 6; j++ {
               w := 3 * pow10[j]
               if w > cap {
                   continue
               }
               cnt := 3 * k
               take := 1
               for cnt > 0 {
                   c := take
                   if c > cnt {
                       c = cnt
                   }
                   weight := c * w
                   profit := int64(c) * F[j]
                   for x := cap; x >= weight; x-- {
                       if dp[x-weight] != inf {
                           v := dp[x-weight] + profit
                           if v > dp[x] {
                               dp[x] = v
                           }
                       }
                   }
                   cnt -= c
                   take <<= 1
               }
           }
           // result is max profit over any used weight <= n
           var res int64
           for w := 0; w <= cap; w++ {
               if dp[w] > res {
                   res = dp[w]
               }
           }
           fmt.Fprintln(os.Stdout, res)
       } else {
           // for small k (<5), use DP over digit positions with carry to match exact n
           // extract digits of n
           nd := [6]int{}
           tmp := n
           for j := 0; j < 6; j++ {
               nd[j] = tmp % 10
               tmp /= 10
           }
           // dpCurr maps carry->max profit
           dpCurr := make(map[int]int64)
           dpCurr[0] = 0
           // positions
           for j := 0; j < 6; j++ {
               dpNext := make(map[int]int64)
               fj := F[j]
               // max sum of digits at pos j
               maxsj := 9 * k
               for cPrev, pPrev := range dpCurr {
                   // s_j ≡ nd[j] - cPrev (mod 10)
                   rem := nd[j] - cPrev
                   rem %= 10
                   if rem < 0 {
                       rem += 10
                   }
                   // try s_j = rem + 10*t
                   for sj := rem; sj <= maxsj; sj += 10 {
                       // next carry
                       t := (sj + cPrev - nd[j]) / 10
                       if t < 0 {
                           continue
                       }
                       // profit add = floor(sj/3)*F[j]
                       add := int64(sj/3) * fj
                       v := pPrev + add
                       if prev, ok := dpNext[t]; !ok || v > prev {
                           dpNext[t] = v
                       }
                   }
               }
               dpCurr = dpNext
           }
           // final carry must be zero
           res := dpCurr[0]
           if res < 0 {
               res = 0
           }
           fmt.Fprintln(os.Stdout, res)
       }
   }
}
