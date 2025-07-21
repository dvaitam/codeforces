package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int64
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   // Precompute factorials up to 20, cap to large value
   const capVal = 1e18
   fact := make([]int64, 21)
   fact[0] = 1
   for i := 1; i <= 20; i++ {
       fact[i] = fact[i-1] * int64(i)
       if fact[i] > capVal {
           fact[i] = capVal
       }
   }
   // Find minimal rem such that rem! >= k
   remMin := -1
   for i := 0; i <= 20; i++ {
       if fact[i] >= k {
           remMin = i
           break
       }
   }
   // If k > n! (when n <=20), remMin will be > n or remain -1
   if remMin < 0 || int64(remMin) > n {
       fmt.Println(-1)
       return
   }
   m := remMin
   prefixLen := n - int64(m)
   // Count lucky positions in prefix where a_i = i
   ans := 0
   for i := int64(1); i <= prefixLen; i++ {
       if isLucky(i) {
           ans++
       }
   }
   // Build tail permutation of size m from values [prefixLen+1 .. n]
   k-- // zero-based index
   // Initialize remaining values
   rem := make([]int64, m)
   for i := 0; i < m; i++ {
       rem[i] = prefixLen + int64(i) + 1
   }
   // Generate k-th permutation of rem
   for i := 0; i < m; i++ {
       f := fact[m-1-i]
       idx := int(k / f)
       val := rem[idx]
       pos := prefixLen + int64(i) + 1
       if isLucky(pos) && isLucky(val) {
           ans++
       }
       // Remove used element
       rem = append(rem[:idx], rem[idx+1:]...)
       k %= f
   }
   fmt.Println(ans)
}

// isLucky returns true if x consists only of digits 4 and 7
func isLucky(x int64) bool {
   if x <= 0 {
       return false
   }
   for x > 0 {
       d := x % 10
       if d != 4 && d != 7 {
           return false
       }
       x /= 10
   }
   return true
}
