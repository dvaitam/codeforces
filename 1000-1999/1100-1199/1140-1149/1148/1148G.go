package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   fmt.Fscan(reader, &n, &k)
   a := make([]int, n)
   maxa := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] > maxa {
           maxa = a[i]
       }
   }
   // sieve primes up to sqrt(maxa)
   lim := intSqrt(maxa)
   sieve := make([]bool, lim+1)
   primes := make([]int, 0)
   for i := 2; i <= lim; i++ {
       if !sieve[i] {
           primes = append(primes, i)
           for j := i * i; j <= lim; j += i {
               sieve[j] = true
           }
       }
   }
   // factor numbers and count prime occurrences
   pf := make([][]int, n)
   cnt2 := make(map[int]int)
   for i := 0; i < n; i++ {
       x := a[i]
       for _, p := range primes {
           if p*p > x {
               break
           }
           if x%p == 0 {
               pf[i] = append(pf[i], p)
               cnt2[p]++
               for x%p == 0 {
                   x /= p
               }
           }
       }
       if x > 1 {
           pf[i] = append(pf[i], x)
           cnt2[x]++
       }
   }
   // inclusion-exclusion counts for all subset products
   cnt := make(map[int]int)
   for i := 0; i < n; i++ {
       m := len(pf[i])
       for mask := 1; mask < (1 << m); mask++ {
           v := 1
           sgn := -1
           for j := 0; j < m; j++ {
               if (mask>>j)&1 == 1 {
                   v *= pf[i][j]
                   sgn = -sgn
               }
           }
           cnt[v] += sgn
       }
   }
   // find element with enough coprime partners
   for i := 0; i < n; i++ {
       m := len(pf[i])
       f := 0
       for mask := 1; mask < (1 << m); mask++ {
           v := 1
           for j := 0; j < m; j++ {
               if (mask>>j)&1 == 1 {
                   v *= pf[i][j]
               }
           }
           f += cnt[v]
       }
       // f is count of numbers (including self) sharing a prime with a[i]
       if n-f+1 >= k {
           // output this index and k-1 coprime indices
           writer.WriteString(fmt.Sprintf("%d", i+1))
           need := k - 1
           for j := 0; j < n && need > 0; j++ {
               if gcd(a[i], a[j]) == 1 {
                   writer.WriteString(fmt.Sprintf(" %d", j+1))
                   need--
               }
           }
           writer.WriteByte('\n')
           return
       }
   }
   // fallback: find smallest prime dividing at least k numbers
   bestP := 0
   for p, c := range cnt2 {
       if c >= k && (bestP == 0 || p < bestP) {
           bestP = p
       }
   }
   if bestP > 0 {
       need := k
       out := make([]int, 0, k)
       for i := 0; i < n && need > 0; i++ {
           if a[i]%bestP == 0 {
               out = append(out, i+1)
               need--
           }
       }
       for idx, v := range out {
           if idx > 0 {
               writer.WriteByte(' ')
           }
           writer.WriteString(fmt.Sprintf("%d", v))
       }
       writer.WriteByte('\n')
   }
}

func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func intSqrt(x int) int {
   r := int(math.Sqrt(float64(x)))
   for (r+1)*(r+1) <= x {
       r++
   }
   for r*r > x {
       r--
   }
   return r
}
