package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   xs := make([]int, n)
   maxX := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &xs[i])
       if xs[i] > maxX {
           maxX = xs[i]
       }
   }
   // count occurrences
   cnt := make([]int32, maxX+1)
   for _, v := range xs {
       cnt[v]++
   }
   // sieve primes up to maxX
   N := maxX
   isPrime := make([]bool, N+1)
   for i := 2; i <= N; i++ {
       isPrime[i] = true
   }
   primes := make([]int, 0, N/10)
   for i := 2; i*i <= N; i++ {
       if isPrime[i] {
           for j := i * i; j <= N; j += i {
               isPrime[j] = false
           }
       }
   }
   for i := 2; i <= N; i++ {
       if isPrime[i] {
           primes = append(primes, i)
       }
   }
   // compute f(p) for each prime
   mP := len(primes)
   f := make([]int32, mP)
   for idx, p := range primes {
       var sum int32
       for j := p; j <= N; j += p {
           sum += cnt[j]
       }
       f[idx] = sum
   }
   // prefix sums
   pre := make([]int64, mP)
   if mP > 0 {
       pre[0] = int64(f[0])
       for i := 1; i < mP; i++ {
           pre[i] = pre[i-1] + int64(f[i])
       }
   }
   // process queries
   var q int
   fmt.Fscan(reader, &q)
   for i := 0; i < q; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       if l > N {
           fmt.Fprintln(writer, 0)
           continue
       }
       if r > N {
           r = N
       }
       // find first prime >= l
       li := sort.Search(mP, func(i int) bool { return primes[i] >= l })
       // find last prime <= r: find first >r, then -1
       ri := sort.Search(mP, func(i int) bool { return primes[i] > r }) - 1
       if li < 0 || li >= mP || ri < li {
           fmt.Fprintln(writer, 0)
       } else {
           var ans int64
           if li > 0 {
               ans = pre[ri] - pre[li-1]
           } else {
               ans = pre[ri]
           }
           fmt.Fprintln(writer, ans)
       }
   }
}
