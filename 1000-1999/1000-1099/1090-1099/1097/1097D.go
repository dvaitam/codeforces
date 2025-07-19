package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const mod = 1000000007

// mul computes (a * b) % mod safely
func mul(a, b, m int64) int64 {
   var res int64
   a %= m
   b %= m
   for b > 0 {
       if b&1 == 1 {
           res += a
           if res >= m {
               res -= m
           }
       }
       a += a
       if a >= m {
           a -= m
       }
       b >>= 1
   }
   return res
}

func powMod(a, b, m int64) int64 {
   var res int64 = 1
   a %= m
   for b > 0 {
       if b&1 == 1 {
           res = mul(res, a, m)
       }
       a = mul(a, a, m)
       b >>= 1
   }
   return res
}

var witnesses = []int64{2, 325, 9375, 28178, 450775, 9780504, 1795265022}

func isPrime(n int64) bool {
   if n < 2 {
       return false
   }
   for _, p := range []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37} {
       if n%p == 0 {
           return n == p
       }
   }
   d := n - 1
   s := 0
   for d&1 == 0 {
       d >>= 1
       s++
   }
   for _, a := range witnesses {
       if a%n == 0 {
           continue
       }
       x := powMod(a, d, n)
       if x == 1 || x == n-1 {
           continue
       }
       composite := true
       for i := 1; i < s; i++ {
           x = mul(x, x, n)
           if x == n-1 {
               composite = false
               break
           }
       }
       if composite {
           return false
       }
   }
   return true
}

func rho(n int64, res *[]int64) {
   if n <= 1 {
       return
   }
   if isPrime(n) {
       *res = append(*res, n)
       return
   }
   var d int64
   for c := int64(1); ; c++ {
       // f(x) = x*x + c mod n
       f := func(x int64) int64 {
           y := mul(x, x, n) + c
           if y >= n {
               y -= n
           }
           return y
       }
       x, y, g := int64(2), int64(2), int64(1)
       for g == 1 {
           x = f(x)
           y = f(f(y))
           diff := x - y
           if diff < 0 {
               diff = -diff
           }
           // gcd
           a, b := diff, n
           for b != 0 {
               a, b = b, a%b
           }
           g = a
       }
       if g > 1 && g < n {
           d = g
           break
       }
   }
   rho(d, res)
   rho(n/d, res)
}

func factorize(n int64) []int64 {
   var smalls []int64
   for i := int64(2); i*i <= n && i <= 100; i++ {
       for n%i == 0 {
           smalls = append(smalls, i)
           n /= i
       }
   }
   if n > 1 {
       rho(n, &smalls)
   }
   sort.Slice(smalls, func(i, j int) bool { return smalls[i] < smalls[j] })
   return smalls
}

func modInv(a int64) int64 {
   return powMod(a, mod-2, mod)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int64
   var k int
   fmt.Fscan(reader, &n, &k)
   primes := factorize(n)
   // precompute inverses
   maxCnt := 0
   // count max multiplicity
   for i := 0; i < len(primes); {
       j := i + 1
       for j < len(primes) && primes[j] == primes[i] {
           j++
       }
       cnt := j - i
       if cnt > maxCnt {
           maxCnt = cnt
       }
       i = j
   }
   preInv := make([]int64, maxCnt+2)
   for i := 1; i < len(preInv); i++ {
       preInv[i] = modInv(int64(i))
   }
   preInv[0] = 1

   result := int64(1)
   // process each prime group
   for i := 0; i < len(primes); {
       p := primes[i]
       j := i + 1
       for j < len(primes) && primes[j] == p {
           j++
       }
       cnt := j - i
       // dp distribution
       dp := make([]int64, cnt+1)
       dp[cnt] = 1
       for rep := 0; rep < k; rep++ {
           var suf int64
           for x := cnt; x >= 0; x-- {
               suf = (suf + dp[x]*preInv[x+1]) % mod
               dp[x] = suf
           }
       }
       // compute contribution
       var total int64
       var mulP int64 = 1
       for x := 0; x <= cnt; x++ {
           total = (total + mulP*dp[x]) % mod
           mulP = mul(mulP, p%mod, mod)
       }
       result = result * total % mod
       i = j
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, result)
}
