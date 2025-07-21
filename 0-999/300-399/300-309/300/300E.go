package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var k int
   if _, err := fmt.Fscan(reader, &k); err != nil {
       return
   }
   a := make([]int, k)
   maxa := 0
   var sumA int64
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] > maxa {
           maxa = a[i]
       }
       sumA += int64(a[i])
   }
   // frequency of each a_i
   freq := make([]int, maxa+2)
   for _, v := range a {
       freq[v]++
   }
   // suffix counts: number of a_i >= v
   suf := make([]int, maxa+2)
   for v := maxa; v >= 0; v-- {
       suf[v] = freq[v] + suf[v+1]
   }
   // sieve primes up to maxa
   isComposite := make([]bool, maxa+1)
   primes := make([]int, 0, maxa/10)
   for i := 2; i <= maxa; i++ {
       if !isComposite[i] {
           primes = append(primes, i)
           if i <= maxa/i {
               for j := i * i; j <= maxa; j += i {
                   isComposite[j] = true
               }
           }
       }
   }
   // compute required prime exponents B_p
   m := len(primes)
   B := make([]int64, m)
   for idx, p := range primes {
       var bp int64
       pPow := p
       for pPow <= maxa {
           var fx int64
           for j := pPow; j <= maxa; j += pPow {
               fx += int64(suf[j])
           }
           bp += fx
           if pPow > maxa/p {
               break
           }
           pPow *= p
       }
       B[idx] = bp
   }
   // binary search minimal n in [1, sumA]
   lo, hi := int64(1), sumA
   var ans int64 = sumA
   for lo <= hi {
       mid := (lo + hi) / 2
       ok := true
       for idx, p := range primes {
           need := B[idx]
           if need == 0 {
               continue
           }
           var have int64
           n := mid
           for n > 0 && have < need {
               n /= int64(p)
               have += n
           }
           if have < need {
               ok = false
               break
           }
       }
       if ok {
           ans = mid
           hi = mid - 1
       } else {
           lo = mid + 1
       }
   }
   fmt.Fprintln(writer, ans)
}
