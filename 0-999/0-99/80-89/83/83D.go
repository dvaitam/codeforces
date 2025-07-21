package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func isPrime(n int64) bool {
   if n < 2 {
       return false
   }
   smallPrimes := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31}
   for _, p := range smallPrimes {
       if n == p {
           return true
       }
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
   // bases for testing deterministic for n < 2^32
   bases := []int64{2, 7, 61}
   for _, a := range bases {
       if a >= n {
           continue
       }
       x := modPow(a, d, n)
       if x == 1 || x == n-1 {
           continue
       }
       composite := true
       for r := 1; r < s; r++ {
           x = modMul(x, x, n)
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

func modMul(a, b, m int64) int64 {
   return (a * b) % m
}

func modPow(a, e, m int64) int64 {
   res := int64(1)
   a %= m
   for e > 0 {
       if e&1 != 0 {
           res = modMul(res, a, m)
       }
       a = modMul(a, a, m)
       e >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, b, k int64
   if _, err := fmt.Fscan(reader, &a, &b, &k); err != nil {
       return
   }
   // if k composite, no numbers
   if !isPrime(k) {
       fmt.Println(0)
       return
   }
   // compute m range [L,R]
   L := (a + k - 1) / k
   R := b / k
   if L > R {
       fmt.Println(0)
       return
   }
   // if k > R, only m==1 possible
   if k > R {
       cnt := int64(0)
       if L <= 1 && 1 <= R {
           cnt = 1
       }
       fmt.Println(cnt)
       return
   }
   // precompute sieve up to maxS
   maxS := int(math.Sqrt(float64(b)/float64(k))) + 2
   // ensure at least up to cbrt(b)
   cb := int(math.Cbrt(float64(b))) + 2
   if cb > maxS {
       maxS = cb
   }
   // also ensure up to sqrt(R) for brute
   tmp := int(math.Sqrt(float64(R))) + 2
   if tmp > maxS {
       maxS = tmp
   }
   sieve := make([]bool, maxS+1)
   primes := make([]int, 0, maxS/10)
   for i := 2; i <= maxS; i++ {
       if !sieve[i] {
           primes = append(primes, i)
           for j := i * i; j <= maxS; j += i {
               sieve[j] = true
           }
       }
   }
   // cubic threshold
   cbrtB := int64(math.Cbrt(float64(b)))
   for (cbrtB+1)*(cbrtB+1)*(cbrtB+1) <= b {
       cbrtB++
   }
   for cbrtB*cbrtB*cbrtB > b {
       cbrtB--
   }
   var ans int64
   // brute if k > cbrt(b)
   if k > cbrtB {
       for m := L; m <= R; m++ {
           if m == 1 {
               ans++
               continue
           }
           if m < k {
               continue
           }
           ok := true
           mm := m
           for _, p := range primes {
               p64 := int64(p)
               if p64 >= k || p64*p64 > mm {
                   break
               }
               if mm%p64 == 0 {
                   ok = false
                   break
               }
           }
           if ok {
               ans++
           }
       }
       fmt.Println(ans)
       return
   }
   // segmented sieve for primes < k
   primesBelow := make([]int64, 0)
   for _, p := range primes {
       if int64(p) < k {
           primesBelow = append(primesBelow, int64(p))
       }
   }
   const blockSize = 1000000
   for start := L; start <= R; start += blockSize {
       end := start + blockSize - 1
       if end > R {
           end = R
       }
       sz := end - start + 1
       mark := make([]bool, sz)
       for _, p := range primesBelow {
           // find first multiple of p in [start,end]
           first := (start + p - 1) / p * p
           for x := first; x <= end; x += p {
               mark[x-start] = true
           }
       }
       for i := int64(0); i < sz; i++ {
           if !mark[i] {
               m := start + i
               if m == 1 || m >= k {
                   ans++
               }
           }
       }
   }
   fmt.Println(ans)
}
