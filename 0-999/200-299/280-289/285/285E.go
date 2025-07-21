package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func add(a, b int) int {
   a += b
   if a >= MOD {
       a -= MOD
   }
   return a
}

func sub(a, b int) int {
   a -= b
   if a < 0 {
       a += MOD
   }
   return a
}

func mul(a, b int) int {
   return int((int64(a) * int64(b)) % MOD)
}

func powmod(a, e int) int {
   res := 1
   x := a
   for e > 0 {
       if e&1 == 1 {
           res = mul(res, x)
       }
       x = mul(x, x)
       e >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   // dp over positions to count valid event subsets of size s
   // dp[s][prevA][prev2A]
   dpCurr := make([][2][2]int, n+1)
   dpCurr[0][0][0] = 1
   for pos := 1; pos <= n; pos++ {
       dpNext := make([][2][2]int, n+1)
       for s := 0; s <= n; s++ {
           for prevA := 0; prevA < 2; prevA++ {
               for prev2A := 0; prev2A < 2; prev2A++ {
                   v := dpCurr[s][prevA][prev2A]
                   if v == 0 {
                       continue
                   }
                   // no event at pos
                   // new prevA = 0, new prev2A = prevA
                   dpNext[s][0][prevA] = add(dpNext[s][0][prevA], v)
                   // event A_pos if pos<n
                   if pos < n {
                       // new prevA=1, new prev2A=prevA
                       dpNext[s+1][1][prevA] = add(dpNext[s+1][1][prevA], v)
                   }
                   // event B_pos if pos>1 and no A_{pos-2} (prev2A==0)
                   if pos > 1 && prev2A == 0 {
                       // new prevA=0, new prev2A=prevA
                       dpNext[s+1][0][prevA] = add(dpNext[s+1][0][prevA], v)
                   }
               }
           }
       }
       dpCurr = dpNext
   }
   // collect dp_events
   dpEvents := make([]int, n+1)
   for s := 0; s <= n; s++ {
       total := 0
       for prevA := 0; prevA < 2; prevA++ {
           for prev2A := 0; prev2A < 2; prev2A++ {
               total = add(total, dpCurr[s][prevA][prev2A])
           }
       }
       dpEvents[s] = total
   }
   // factorials and invfacts
   fact := make([]int, n+1)
   invFact := make([]int, n+1)
   fact[0] = 1
   for i := 1; i <= n; i++ {
       fact[i] = mul(fact[i-1], i)
   }
   invFact[n] = powmod(fact[n], MOD-2)
   for i := n; i > 0; i-- {
       invFact[i-1] = mul(invFact[i], i)
   }
   // compute result
   res := 0
   for s := k; s <= n; s++ {
       ds := dpEvents[s]
       if ds == 0 {
           continue
       }
       // C(s, k)
       comb := mul(fact[s], mul(invFact[k], invFact[s-k]))
       term := mul(ds, comb)
       term = mul(term, fact[n-s])
       if (s-k)&1 == 1 {
           res = sub(res, term)
       } else {
           res = add(res, term)
       }
   }
   // print
   fmt.Println(res)
}
