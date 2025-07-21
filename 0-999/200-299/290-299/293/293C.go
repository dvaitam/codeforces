package main

import (
   "bufio"
   "fmt"
   "math/big"
   "math/rand"
   "os"
   "sort"
   "time"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int64
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // 3*(a+b)(b+c)(c+a) = n
   if n%3 != 0 {
       fmt.Println(0)
       return
   }
   m := n / 3
   // factor m
   fs := make(map[int64]int)
   rand.Seed(time.Now().UnixNano())
   factor(m, fs)
   // generate divisors
   divs := divisorsFromFactors(fs)
   sort.Slice(divs, func(i, j int) bool { return divs[i] < divs[j] })
   var cnt int64
   for _, x := range divs {
       if x > m {
           break
       }
       t1 := m / x
       for _, y := range divs {
           if y > t1 {
               break
           }
           if t1%y != 0 {
               continue
           }
           z := t1 / y
           // check parity
           if (x+y+z)&1 != 0 {
               continue
           }
           // compute a,b,c
           a := (x + z - y) / 2
           b := (x + y - z) / 2
           c := (y + z - x) / 2
           if a > 0 && b > 0 && c > 0 {
               cnt++
           }
       }
   }
   fmt.Println(cnt)
}

// divisorsFromFactors returns all positive divisors of the number described by prime factors fs
func divisorsFromFactors(fs map[int64]int) []int64 {
   divs := []int64{1}
   for p, e := range fs {
       sz := len(divs)
       var mul int64 = 1
       for i := 1; i <= e; i++ {
           mul *= p
           for j := 0; j < sz; j++ {
               divs = append(divs, divs[j]*mul)
           }
       }
   }
   return divs
}

// factor performs Pollard's Rho factorization, filling map fs with prime factors
func factor(n int64, fs map[int64]int) {
   if n <= 1 {
       return
   }
   if isPrime(n) {
       fs[n]++
   } else {
       d := rho(n)
       factor(d, fs)
       factor(n/d, fs)
   }
}

// isPrime tests primality using Miller-Rabin
func isPrime(n int64) bool {
   if n < 2 {
       return false
   }
   smallPrimes := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
   for _, p := range smallPrimes {
       if n == p {
           return true
       }
       if n%p == 0 {
           return false
       }
   }
   d := n - 1
   s := 0
   for d&1 == 0 {
       d >>= 1
       s++
   }
   // test bases
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

// rho finds a non-trivial divisor of n
func rho(n int64) int64 {
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
           d = gcd(abs(x-y), n)
           if d == n {
               break
           }
       }
       if d > 1 && d < n {
           return d
       }
   }
}

func modMul(a, b, mod int64) int64 {
   // compute (a * b) % mod safely
   t := new(big.Int).Mul(big.NewInt(a), big.NewInt(b))
   t.Mod(t, big.NewInt(mod))
   return t.Int64()
}

func modPow(a, d, mod int64) int64 {
   result := int64(1)
   base := a % mod
   for d > 0 {
       if d&1 == 1 {
           result = modMul(result, base, mod)
       }
       base = modMul(base, base, mod)
       d >>= 1
   }
   return result
}

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func abs(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}
