package main

import (
   "bufio"
   "fmt"
   "math"
   "math/big"
   "math/rand"
   "os"
   "sort"
   "time"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var A int64
   fmt.Fscan(in, &A)
   if A == 1 {
       fmt.Println(1)
       return
   }
   // factor A
   rand.Seed(time.Now().UnixNano())
   facA := factor(A)
   // generate divisors
   divs := make([]int64, 0)
   genDivs(facA, 0, 1, &divs)
   // prepare candidates d where d>1 and d-1 is prime power
   type cand struct{ d, p int64 }
   cands := make([]cand, 0)
   for _, d := range divs {
       if d <= 1 {
           continue
       }
       v := d - 1
       if v <= 1 {
           continue
       }
       mp := factor(v)
       if len(mp) == 1 {
           for p := range mp {
               cands = append(cands, cand{d: d, p: p})
           }
       }
   }
   sort.Slice(cands, func(i, j int) bool { return cands[i].d < cands[j].d })
   // dfs to count ways
   used := make(map[int64]bool)
   var dfs func(rem int64, idx int) int64
   dfs = func(rem int64, idx int) int64 {
       if rem == 1 {
           return 1
       }
       var cnt int64
       for i := idx; i < len(cands); i++ {
           d := cands[i].d
           if d > rem {
               break
           }
           if rem%d != 0 {
               continue
           }
           p := cands[i].p
           if used[p] {
               continue
           }
           used[p] = true
           cnt += dfs(rem/d, i+1)
           used[p] = false
       }
       return cnt
   }
   res := dfs(A, 0)
   fmt.Println(res)
}

// genDivs generates all divisors from prime factors
func genDivs(mp map[int64]int, i int, cur int64, res *[]int64) {
   if i == len(mp) {
       *res = append(*res, cur)
       return
   }
   // convert map to slice once
   keys := make([]int64, 0, len(mp))
   for k := range mp {
       keys = append(keys, k)
   }
   sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
   // recursion over exponents
   p := keys[i]
   e := mp[p]
   val := int64(1)
   for k := 0; k <= e; k++ {
       genDivs(mp, i+1, cur*val, res)
       val *= p
   }
}

// Miller-Rabin and Pollard Rho
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
   for d&1 == 0 {
       d >>= 1
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
       composite := true
       for r := 1; r < s; r++ {
           x = mulMod(x, x, n)
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

func modPow(a, d, mod int64) int64 {
   result := int64(1)
   a %= mod
   for d > 0 {
       if d&1 == 1 {
           result = mulMod(result, a, mod)
       }
       a = mulMod(a, a, mod)
       d >>= 1
   }
   return result
}

func mulMod(a, b, mod int64) int64 {
   return int64(new(big.Int).Mod(new(big.Int).Mul(big.NewInt(a), big.NewInt(b)), big.NewInt(mod)).Int64())
}

func pollardsRho(n int64) int64 {
   if n%2 == 0 {
       return 2
   }
   x := rand.Int63n(n-2) + 2
   y := x
   c := rand.Int63n(n-1) + 1
   d := int64(1)
   for d == 1 {
       x = (mulMod(x, x, n) + c) % n
       y = (mulMod(y, y, n) + c) % n
       y = (mulMod(y, y, n) + c) % n
       d = gcd(abs(x-y), n)
       if d == n {
           return pollardsRho(n)
       }
   }
   return d
}

func factor(n int64) map[int64]int {
   res := make(map[int64]int)
   var f func(int64)
   f = func(n int64) {
       if n == 1 {
           return
       }
       if isPrime(n) {
           res[n]++
           return
       }
       d := pollardsRho(n)
       f(d)
       f(n / d)
   }
   f(n)
   return res
}

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func abs(a int64) int64 {
   if a < 0 {
       return -a
   }
   return a
}
