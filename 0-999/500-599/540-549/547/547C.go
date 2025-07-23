package main

import (
   "bufio"
   "fmt"
   "os"
)

// pair holds a divisor product and its inclusion-exclusion sign
type pair struct { prod, sign int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, q int
   fmt.Fscan(reader, &n, &q)
   a := make([]int, n)
   maxA := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] > maxA {
           maxA = a[i]
       }
   }
   // sieve smallest prime factors
   spf := make([]int, maxA+1)
   for i := 2; i <= maxA; i++ {
       if spf[i] == 0 {
           for j := i; j <= maxA; j += i {
               if spf[j] == 0 {
                   spf[j] = i
               }
           }
       }
   }
   // precompute subsets of prime factors for each a[i]
   subsList := make([][]pair, n)
   for i := 0; i < n; i++ {
       x := a[i]
       // collect distinct primes
       var primes []int
       for x > 1 {
           p := spf[x]
           primes = append(primes, p)
           for x%p == 0 {
               x /= p
           }
       }
       k := len(primes)
       if k > 0 {
           // generate non-empty subsets
           var subs []pair
           // mask from 1 to (1<<k)-1
           for mask := 1; mask < (1 << k); mask++ {
               prod := 1
               bits := 0
               for j := 0; j < k; j++ {
                   if mask&(1<<j) != 0 {
                       prod *= primes[j]
                       bits++
                   }
               }
               sign := 1
               if bits%2 == 0 {
                   sign = -1
               }
               subs = append(subs, pair{prod, sign})
           }
           subsList[i] = subs
       }
   }
   cnt := make([]int, maxA+1)
   inShelf := make([]bool, n)
   currSize := 0
   var ans int64

   // process queries
   for qi := 0; qi < q; qi++ {
       var x int
       fmt.Fscan(reader, &x)
       idx := x - 1
       subs := subsList[idx]
       // compute number of non-coprime elements via inclusion-exclusion
       ie := 0
       for _, p := range subs {
           ie += p.sign * cnt[p.prod]
       }
       // number of coprime with a[idx]
       cp := currSize - ie
       if !inShelf[idx] {
           // add
           ans += int64(cp)
           // update counts
           for _, p := range subs {
               cnt[p.prod]++
           }
           currSize++
           inShelf[idx] = true
       } else {
           // remove
           // special case for a[i]==1 (no primes): avoid counting self
           if len(subs) == 0 {
               cp--
           }
           ans -= int64(cp)
           // update counts
           for _, p := range subs {
               cnt[p.prod]--
           }
           currSize--
           inShelf[idx] = false
       }
       fmt.Fprintln(writer, ans)
   }
}
