package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "sort"
   "time"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   if n == 1 {
       fmt.Println(a[0])
       return
   }
   half := (n + 1) / 2
   rand.Seed(time.Now().UnixNano())
   best := int64(1)
   samples := 10
   for t := 0; t < samples; t++ {
       x := a[rand.Intn(n)]
       pf := factor(x)
       ds := divisors(pf)
       sort.Slice(ds, func(i, j int) bool { return ds[i] > ds[j] })
       for _, d := range ds {
           if d <= best {
               break
           }
           cnt := 0
           for _, v := range a {
               if v%d == 0 {
                   cnt++
                   if cnt >= half {
                       best = d
                       break
                   }
               }
           }
       }
   }
   fmt.Println(best)
}

// factor returns prime factorization of n as map prime->exponent
func factor(n int64) map[int64]int {
   res := make(map[int64]int)
   var dfs func(int64)
   dfs = func(n int64) {
       if n == 1 {
           return
       }
       if isPrime(n) {
           res[n]++
       } else {
           d := pollardsRho(n)
           dfs(d)
           dfs(n / d)
       }
   }
   dfs(n)
   return res
}

// divisors returns all divisors from prime factor map
func divisors(pf map[int64]int) []int64 {
   ds := []int64{1}
   for p, e := range pf {
       cur := make([]int64, 0, len(ds)*(e+1))
       pow := int64(1)
       for i := 0; i <= e; i++ {
           for _, d := range ds {
               cur = append(cur, d*pow)
           }
           pow *= p
       }
       ds = cur
   }
   return ds
}

// isPrime Miller-Rabin for 64-bit
func isPrime(n int64) bool {
   if n < 2 {
       return false
   }
   smallPrimes := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
   for _, p := range smallPrimes {
       if n%p == 0 {
           return n == p
       }
   }
   d := n - 1
   s := 0
   for d%2 == 0 {
       d /= 2
       s++
   }
   bases := []int64{2, 325, 9375, 28178, 450775, 9780504, 1795265022}
   for _, a := range bases {
       if a%n == 0 {
           continue
       }
       x := modPow(a, d, n)
       if x == 1 || x == n-1 {
           continue
       }
       skip := false
       for r := 1; r < s; r++ {
           x = modMul(x, x, n)
           if x == n-1 {
               skip = true
               break
           }
       }
       if skip {
           continue
       }
       return false
   }
   return true
}

// modMul returns (a * b) % m avoiding overflow
func modMul(a, b, m int64) int64 {
   var res int64
   a %= m
   for b > 0 {
       if b&1 == 1 {
           res = (res + a) % m
       }
       a = (a << 1) % m
       b >>= 1
   }
   return res
}

func modPow(a, d, m int64) int64 {
   result := int64(1)
   a %= m
   for d > 0 {
       if d&1 == 1 {
           result = modMul(result, a, m)
       }
       a = modMul(a, a, m)
       d >>= 1
   }
   return result
}

// pollardsRho finds a divisor of n
func pollardsRho(n int64) int64 {
   if n%2 == 0 {
       return 2
   }
   for {
       x := rand.Int63n(n-2) + 2
       y := x
       c := rand.Int63n(n-1) + 1
       d := int64(1)
       for d == 1 {
           x = (modMul(x, x, n) + c) % n
           y = (modMul(y, y, n) + c) % n
           y = (modMul(y, y, n) + c) % n
           diff := x - y
           if diff < 0 {
               diff = -diff
           }
           d = gcd(diff, n)
           if d == n {
               break
           }
       }
       if d > 1 && d < n {
           return d
       }
   }
}

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}
