package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   const MOD = 1 << 30
   in := bufio.NewReader(os.Stdin)
   var a, b, c int
   fmt.Fscan(in, &a, &b, &c)
   ab := a * b
   // cnt[t] = number of pairs (i,j) with i*j = t
   cnt := make([]int32, ab+1)
   for i := 1; i <= a; i++ {
       for j := 1; j <= b; j++ {
           cnt[i*j]++
       }
   }
   // D[u] = number of pairs (i,j) such that u divides i*j
   // which is sum_{k: u*k <= ab} cnt[u*k]
   D := make([]int32, ab+1)
   for u := 1; u <= ab; u++ {
       var sum int32
       for v := u; v <= ab; v += u {
           sum += cnt[v]
       }
       D[u] = sum
   }
   // compute mobius mu up to c
   mu := make([]int8, c+1)
   isPrime := make([]bool, c+1)
   primes := make([]int, 0, c/10)
   if c >= 1 {
       mu[1] = 1
   }
   for i := 2; i <= c; i++ {
       isPrime[i] = true
   }
   for i := 2; i <= c; i++ {
       if isPrime[i] {
           primes = append(primes, i)
           mu[i] = -1
       }
       for _, p := range primes {
           v := i * p
           if v > c {
               break
           }
           isPrime[v] = false
           if i%p == 0 {
               mu[v] = 0
               break
           }
           mu[v] = -mu[i]
       }
   }
   // H[m] = sum_{k=1..m} floor(m/k)
   H := make([]int32, c+1)
   for m := 1; m <= c; m++ {
       var s int32
       for k := 1; k <= m; k++ {
           s += int32(m / k)
       }
       H[m] = s
   }
   // T[u] = sum_{d|u, d<=c} mu[d] * H[c/d]
   T := make([]int64, ab+1)
   for d := 1; d <= c; d++ {
       md := mu[d]
       if md == 0 {
           continue
       }
       h := int64(H[c/d]) * int64(md)
       for u := d; u <= ab; u += d {
           T[u] += h
       }
   }
   // accumulate answer
   var ans int64
   for u := 1; u <= ab; u++ {
       if D[u] != 0 && T[u] != 0 {
           ans += int64(D[u]) * T[u]
           // avoid overflow
           if ans > (1<<61) || ans < -(1<<61) {
               ans %= MOD
           }
       }
   }
   ans %= MOD
   if ans < 0 {
       ans += MOD
   }
   fmt.Println(ans)
}
