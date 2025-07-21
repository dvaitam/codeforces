package main

import (
   "bufio"
   "fmt"
   "os"
)

// solve one test case for given n,m
func solveCase(n, m int64) int64 {
   if n <= 1 {
       return 0
   }
   // try small k up to k_limit; beyond that, N(k)=k+1
   var kLimit int64 = 60
   if m < kLimit {
       kLimit = m
   }
   for k := int64(1); k <= kLimit; k++ {
       // find t_max: largest t in [0..k] such that S(k-1, t-1) <= m
       // where S(i, j) = sum_{u=0..j} C(i, u)
       var sPrev int64 = 0
       var comb int64 = 1 // C(k-1, 0)
       tMax := int64(0)
       // t=0: sPrev=0 <= m
       tMax = 0
       // iterate t from 1..k
       for t := int64(1); t <= k; t++ {
           // compute comb = C(k-1, t-1)
           if t == 1 {
               comb = 1
           } else {
               // update comb to C(k-1, t-1)
               // comb_prev was C(k-1, t-2)
               // comb = comb_prev * ((k-1) - (t-2)) / (t-1)
               num := (k-1) - (t - 2)
               den := t - 1
               // multiply then divide, watch for overflow but values small
               comb = comb * num / den
           }
           // accumulate
           if sPrev + comb > m {
               break
           }
           sPrev += comb
           tMax = t
       }
       // now tMax is largest t with S(k-1, t-1)<=m
       // N = S(k, tMax) = sum_{i=0..tMax} C(k, i)
       var total int64 = 0
       comb = 1 // C(k, 0)
       for i := int64(0); i <= tMax; i++ {
           if i > 0 {
               // update comb = C(k, i)
               num := k - (i - 1)
               den := i
               comb = comb * num / den
           }
           total += comb
           if total >= n {
               break
           }
       }
       if total >= n {
           return k
       }
   }
   // no small k works, so must have k > m and N(k)=k+1 => k+1 >= n => k >= n-1
   return n - 1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n, m int64
       fmt.Fscan(reader, &n, &m)
       res := solveCase(n, m)
       fmt.Fprintln(writer, res)
   }
}
