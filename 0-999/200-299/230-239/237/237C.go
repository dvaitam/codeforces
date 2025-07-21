package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, b, k int
   if _, err := fmt.Fscan(reader, &a, &b, &k); err != nil {
       return
   }
   n := b - a + 1
   // Sieve primes up to b
   isPrime := make([]bool, b+1)
   for i := 2; i <= b; i++ {
       isPrime[i] = true
   }
   for i := 2; i*i <= b; i++ {
       if isPrime[i] {
           for j := i * i; j <= b; j += i {
               isPrime[j] = false
           }
       }
   }
   // Prefix sum of primes
   prefix := make([]int, b+1)
   for i := 1; i <= b; i++ {
       prefix[i] = prefix[i-1]
       if isPrime[i] {
           prefix[i]++
       }
   }
   // Binary search on l
   lo, hi := k, n
   ans := -1
   for lo <= hi {
       mid := (lo + hi) / 2
       if ok(a, b, k, mid, prefix) {
           ans = mid
           hi = mid - 1
       } else {
           lo = mid + 1
       }
   }
   fmt.Println(ans)
}

// ok checks if every subarray of length l in [a,b] has at least k primes
func ok(a, b, k, l int, prefix []int) bool {
   end := b - l + 1
   for x := a; x <= end; x++ {
       // range [x, x+l-1]
       cnt := prefix[x+l-1] - prefix[x-1]
       if cnt < k {
           return false
       }
   }
   return true
}
