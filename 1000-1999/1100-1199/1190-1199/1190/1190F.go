package main

import (
   "bufio"
   "fmt"
   "math"
   "math/big"
   "os"
)

// find prime base q and exponent t such that m = q^t
func factorPrimePower(m uint64) (uint64, int) {
   mBig := new(big.Int).SetUint64(m)
   // try t from high to low
   for t := 60; t >= 2; t-- {
       // approximate t-th root
       r := uint64(math.Round(math.Pow(float64(m), 1.0/float64(t))))
       for d := int64(r) - 1; d <= int64(r)+1; d++ {
           if d <= 1 {
               continue
           }
           cand := new(big.Int).SetInt64(d)
           // cand^t
           pow := new(big.Int).Exp(cand, big.NewInt(int64(t)), nil)
           if pow.Cmp(mBig) == 0 {
               return uint64(d), t
           }
       }
   }
   // m is prime
   return m, 1
}

// factor small uint64 x into distinct prime factors (trial on q-1)
func factorDistinct(x uint64) []uint64 {
   var fs []uint64
   if x%2 == 0 {
       fs = append(fs, 2)
       for x%2 == 0 {
           x /= 2
       }
   }
   for i := uint64(3); i*i <= x; i += 2 {
       if x%i == 0 {
           fs = append(fs, i)
           for x%i == 0 {
               x /= i
           }
       }
   }
   if x > 1 {
       fs = append(fs, x)
   }
   return fs
}

// modular exponentiation for uint64 using big.Int (occasional use)
func modPow(a, e, m uint64) uint64 {
   return new(big.Int).Exp(
       new(big.Int).SetUint64(a),
       new(big.Int).SetUint64(e),
       new(big.Int).SetUint64(m)).Uint64()
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m, p uint64
   if _, err := fmt.Fscan(in, &n, &m, &p); err != nil {
       return
   }
   // factor m = q^t
   q, t := factorPrimePower(m)
   // phi = q^{t-1} * (q-1)
   var phi uint64 = 1
   for i := 0; i < t-1; i++ {
       phi *= q
   }
   phi *= (q - 1)
   // case p divisible by q: only forbidden are x not coprime and x==1
   if p%q == 0 {
       if phi < 1+n {
           fmt.Fprintln(out, -1)
           return
       }
       // enumerate x >=2, x<m, x%q!=0
       cnt := uint64(0)
       for x := uint64(2); x < m && cnt < n; x++ {
           if x%q == 0 {
               continue
           }
           fmt.Fprintln(out, x)
           cnt++
       }
       return
   }
   // gcd(p,m)==1: compute order of p mod m
   // factor phi into distinct primes
   var primes []uint64
   if t > 1 {
       primes = append(primes, q)
   }
   // factors of q-1
   primes = append(primes, factorDistinct(q-1)...)   
   // remove duplicates
   mset := make(map[uint64]bool)
   for _, v := range primes {
       mset[v] = true
   }
   primes = primes[:0]
   for v := range mset {
       primes = append(primes, v)
   }
   // initial order
   ord := phi
   // reduce ord
   for _, f := range primes {
       for ord%f == 0 {
           if modPow(p, ord/f, m) == 1 {
               ord /= f
           } else {
               break
           }
       }
   }
   // available count
   if phi < ord+n {
       fmt.Fprintln(out, -1)
       return
   }
   // prepare big ints
   mBig := new(big.Int).SetUint64(m)
   ordBig := new(big.Int).SetUint64(ord)
   oneBig := big.NewInt(1)
   xBig := new(big.Int)
   // enumerate x
   cnt := uint64(0)
   for x := uint64(1); x < m && cnt < n; x++ {
       if x%q == 0 {
           continue
       }
       xBig.SetUint64(x)
       if new(big.Int).Exp(xBig, ordBig, mBig).Cmp(oneBig) == 0 {
           continue
       }
       fmt.Fprintln(out, x)
       cnt++
   }
}
