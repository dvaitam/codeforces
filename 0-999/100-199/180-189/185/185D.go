package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

func modExp(a, e, mod int64) int64 {
   res := int64(1)
   a %= mod
   for e > 0 {
       if e&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       e >>= 1
   }
   return res
}

// compute order of 2 mod p, p prime
func order2(p int64) int64 {
   phi := p - 1
   // factor phi
   var fac []int64
   x := phi
   for i := int64(2); i*i <= x; i++ {
       if x%i == 0 {
           fac = append(fac, i)
           for x%i == 0 {
               x /= i
           }
       }
   }
   if x > 1 {
       fac = append(fac, x)
   }
   ord := phi
   for _, pr := range fac {
       for ord%pr == 0 {
           if modExp(2, ord/pr, p) == 1 {
               ord /= pr
           } else {
               break
           }
       }
   }
   return ord
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   for ti := 0; ti < t; ti++ {
       var k, l, r, p int64
       fmt.Fscan(in, &k, &l, &r, &p)
       if k%p == 0 {
           // Ai = 1 mod p
           fmt.Fprintln(out, 1)
           continue
       }
       // find ord of 2 mod p
       ord := order2(p)
       if r-l+1 >= ord {
           fmt.Fprintln(out, 0)
           continue
       }
       // compute lcm mod p via big ints for small ranges
       res := big.NewInt(1)
       km := big.NewInt(k)
       two := big.NewInt(2)
       for i := int64(0); i <= r-l; i++ {
           // Ai = k*2^(l+i) +1
           // compute 2^(l+i) as big.Int
           exp := new(big.Int).Exp(two, big.NewInt(l+i), nil)
           ai := new(big.Int).Mul(km, exp)
           ai.Add(ai, big.NewInt(1))
           if new(big.Int).Mod(ai, big.NewInt(p)).Sign() == 0 {
               res.SetInt64(0)
               break
           }
           // gcd(res, ai)
           g := new(big.Int).GCD(nil, nil, res, ai)
           // res = res * ai / g
           res.Mul(res, ai)
           res.Div(res, g)
           // to avoid blowup, we can mod by p but keep value for gcd
           // but for simplicity, no mod
       }
       // print res mod p
       fmt.Fprintln(out, new(big.Int).Mod(res, big.NewInt(p)))
   }
}
