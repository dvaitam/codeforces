package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

// Count decompositions of n into distinct Fibonacci numbers F1=1, F2=2, Fi=Fi-1+Fi-2
var fib []int64
var sumFib []int64
var memo map[key]*big.Int

type key struct{ n int64; k int }

// f(n,k) returns number of ways to write n using fib[0..k]
func f(n int64, k int) *big.Int {
   if n == 0 {
       return big.NewInt(1)
   }
   if k < 0 {
       return big.NewInt(0)
   }
   // prune: if sum of all fib[0..k] < n, impossible
   if sumFib[k] < n {
       return big.NewInt(0)
   }
   key := key{n, k}
   if v, ok := memo[key]; ok {
       return v
   }
   var res *big.Int
   if fib[k] > n {
       res = f(n, k-1)
   } else {
       a := f(n, k-1)
       b := f(n-fib[k], k-1)
       res = new(big.Int).Add(a, b)
   }
   memo[key] = res
   return res
}

func main() {
   // precompute fibs up to >1e18
   fib = make([]int64, 0, 100)
   fib = append(fib, 1, 2)
   for {
       i := len(fib)
       next := fib[i-1] + fib[i-2]
       if next < 0 || next > 1e18 {
           break
       }
       fib = append(fib, next)
   }
   // prefix sums
   sumFib = make([]int64, len(fib))
   for i, v := range fib {
       if i == 0 {
           sumFib[i] = v
       } else {
           sumFib[i] = sumFib[i-1] + v
       }
   }
   memo = make(map[key]*big.Int)

   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var n int64
       fmt.Fscan(in, &n)
       // find largest k with fib[k] <= n
       k := len(fib) - 1
       for k >= 0 && fib[k] > n {
           k--
       }
       ans := f(n, k)
       out.WriteString(ans.String())
       out.WriteByte('\n')
   }
}
