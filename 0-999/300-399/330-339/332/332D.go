package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(reader, &n, &k)
   deg := make([]int, n)
   sumW := make([]int64, n)
   // read upper triangle
   for i := 0; i < n-1; i++ {
       for j := i+1; j < n; j++ {
           var c int64
           fmt.Fscan(reader, &c)
           if c != -1 {
               deg[i]++
               deg[j]++
               sumW[i] += c
               sumW[j] += c
           }
       }
   }
   // precompute factorials up to max(n, max deg)
   maxN := n
   for _, d := range deg {
       if d > maxN {
           maxN = d
       }
   }
   fact := make([]*big.Int, maxN+1)
   fact[0] = big.NewInt(1)
   for i := 1; i <= maxN; i++ {
       fact[i] = new(big.Int).Mul(fact[i-1], big.NewInt(int64(i)))
   }
   // C(n, k)
   denom := binomBig(fact, n, k)
   total := big.NewInt(0)
   tmp := new(big.Int)
   for v := 0; v < n; v++ {
       if deg[v] >= k {
           // C(deg[v]-1, k-1)
           c := binomBig(fact, deg[v]-1, k-1)
           // multiply by sumW[v]
           tmp.Mul(c, big.NewInt(sumW[v]))
           total.Add(total, tmp)
       }
   }
   // compute floor(total/denom)
   result := new(big.Int).Div(total, denom)
   // adjust for negative total to floor
   if total.Sign() < 0 {
       rem := new(big.Int).Mod(total, denom)
       if rem.Sign() != 0 {
           result.Sub(result, big.NewInt(1))
       }
   }
   fmt.Println(result)
}

// binomBig computes C(n, k) = fact[n] / (fact[k] * fact[n-k])
func binomBig(fact []*big.Int, n, k int) *big.Int {
   if k < 0 || k > n {
       return big.NewInt(0)
   }
   res := new(big.Int).Set(fact[n])
   denom := new(big.Int).Mul(fact[k], fact[n-k])
   res.Div(res, denom)
   return res
}
