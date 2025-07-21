package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// fast pow mod
func modPow(a, e, mod int64) int64 {
   res := int64(1)
   a %= mod
   for e > 0 {
       if e&1 == 1 {
           res = (res * a) % mod
       }
       a = (a * a) % mod
       e >>= 1
   }
   return res
}

// factorize n into map prime->exponent
func factorize(n int64) map[int64]int {
   m := make(map[int64]int)
   x := n
   for d := int64(2); d*d <= x; d++ {
       for x%d == 0 {
           m[d]++
           x /= d
       }
   }
   if x > 1 {
       m[x]++
   }
   return m
}

func genDivs(pr map[int64]int) []int64 {
   divs := []int64{1}
   for p, e := range pr {
       sz := len(divs)
       powp := int64(1)
       for i := 1; i <= e; i++ {
           powp *= p
           for j := 0; j < sz; j++ {
               divs = append(divs, divs[j]*powp)
           }
       }
   }
   sort.Slice(divs, func(i, j int) bool { return divs[i] < divs[j] })
   return divs
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   var p int64
   fmt.Fscan(in, &n, &m, &p)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   b := make([]int64, m)
   zero := false
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &b[i])
       b[i] %= p
       if b[i] == 0 {
           zero = true
       }
   }
   // group size p-1
   p1 := p - 1
   // factorize p1
   pr := factorize(p1)
   // include zero if any ai or bj is zero mod p
   for i := 0; i < n; i++ {
       if a[i]%p == 0 {
           zero = true
           break
       }
   }
   // find first non-zero b
   var b0 int64 = 0
   nz := 0
   for j := 0; j < m; j++ {
       if b[j] != 0 {
           b0 = b[j]
           nz = 1
           break
       }
   }
   // if no non-zero b, only identity in multiplicative part
   if nz == 0 {
       // count identity + possible zero
       res := int64(1)
       if zero {
           res++
       }
       fmt.Println(res)
       return
   }
   // build c_j = b_j * inv(b0) for non-zero b_j
   invB0 := modPow(b0, p-2, p)
   cs := make([]int64, 0, m)
   for j := 0; j < m; j++ {
       if b[j] != 0 {
           cs = append(cs, (b[j]*invB0)%p)
       }
   }
   // compute gB = gcd of exponents of cs entries and p1
   gB := int64(1)
   for q, e := range pr {
       powq := int64(1)
       cnt := 0
       for t := 1; t <= e; t++ {
           powq *= q
           exp := p1 / powq
           ok := true
           for _, cj := range cs {
               if modPow(cj, exp, p) != 1 {
                   ok = false
                   break
               }
           }
           if !ok {
               break
           }
           cnt++
       }
       for i := 0; i < cnt; i++ {
           gB *= q
       }
   }
   // compute divisors of p1
   divs := genDivs(pr)
   // map divisor to index
   idx := make(map[int64]int)
   for i, d := range divs {
       idx[d] = i
   }
   // compute phi for each divisor
   phi := make([]int64, len(divs))
   for i, d := range divs {
       // phi(d)
       res := d
       // factorize d using pr
       tmp := d
       for q := range pr {
           if tmp%q == 0 {
               res = res / q * (q - 1)
               for tmp%q == 0 {
                   tmp /= q
               }
           }
       }
       phi[i] = res
   }
   // compute D = set of d_i
   Dmap := make(map[int64]bool)
   // factorize gB to get prime exponents
   gbpr := factorize(gB)
   for i := 0; i < n; i++ {
       ai := a[i] % p
       var di int64
       if ai == 0 {
           // only identity in multiplicative part
           di = p1
       } else {
           // y0 exponent corresponds to e_ai + B0
           y0 := (ai * b0) % p
           di = 1
           // for each prime in gB, find exponent in gcd
           for q, e := range gbpr {
               powq := int64(1)
               tmax := 0
               for t := 1; t <= e; t++ {
                   powq *= q
                   if modPow(y0, p1/powq, p) == 1 {
                       tmax = t
                   } else {
                       break
                   }
               }
               for k := 0; k < tmax; k++ {
                   di *= q
               }
           }
       }
       Dmap[di] = true
   }
   // collect minimal divisors? Actually for marking t where any d divides t
   // flag for each divisor t of p1
   flag := make([]bool, len(divs))
   for d := range Dmap {
       // for each multiple t of d
       for _, t := range divs {
           if t%d == 0 {
               flag[idx[t]] = true
           }
       }
   }
   // sum phi(p1/t) for flagged t
   ans := int64(0)
   for i, t := range divs {
       if flag[i] {
           // k = p1/t
           k := p1 / t
           // phi(k)
           // find index of k in divs
           j := idx[k]
           ans += phi[j]
       }
   }
   if zero {
       ans++
   }
   fmt.Println(ans)
}
