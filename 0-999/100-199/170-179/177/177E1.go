package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var c int64
   if _, err := fmt.Fscan(in, &n, &c); err != nil {
       return
   }
   a := make([]int64, n)
   b := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i], &b[i])
   }
   // Total days = n + sum floor(a[i]*x / b[i]) == c
   // Let T = c - n; need sum floors = T
   T := c - int64(n)
   if T < 0 {
       fmt.Println(0)
       return
   }
   // helper: compute sum floor(a[i]*x/b[i]), capped at T+1
   sumFloor := func(x int64) int64 {
       var sum int64
       for i := 0; i < n; i++ {
           ai := a[i]
           bi := b[i]
           if ai == 0 {
               continue
           }
           // compute term = ai*x/bi
           term := ai * x / bi
           if term > T {
               return T + 1
           }
           sum += term
           if sum > T {
               return T + 1
           }
       }
       return sum
   }
   // find first x>=1 such that sumFloor(x) >= T
   // doubling to find hi
   lo := int64(1)
   hi := int64(1)
   for sumFloor(hi) < T {
       hi <<= 1
       if hi > 2000000000 {
           break
       }
   }
   // if we broke due to hi limit and sumFloor still < T => no solution
   if sumFloor(hi) < T {
       // only possible infinite if T==0 and all a[i]==0
       if T == 0 {
           // check if all a are zero
           allZero := true
           for i := 0; i < n; i++ {
               if a[i] != 0 {
                   allZero = false
                   break
               }
           }
           if allZero {
               fmt.Println(-1)
               return
           }
           // finite but no x gives sum>=0? impossible
       }
       fmt.Println(0)
       return
   }
   // binary search for first >= T
   for lo < hi {
       mid := lo + (hi-lo)/2
       if sumFloor(mid) < T {
           lo = mid + 1
       } else {
           hi = mid
       }
   }
   x0 := lo
   if sumFloor(x0) != T {
       fmt.Println(0)
       return
   }
   // find first x such that sumFloor(x) > T
   lo = x0
   // hi already >= x0 and sumFloor(hi)>=T, but we need >T
   // ensure hi is such that sumFloor(hi)>T
   if sumFloor(hi) == T {
       for sumFloor(hi) <= T {
           hi <<= 1
           if hi > 2000000000 {
               break
           }
       }
   }
   if sumFloor(hi) <= T {
       // infinite
       fmt.Println(-1)
       return
   }
   // binary search for first >T in [lo, hi]
   for lo < hi {
       mid := lo + (hi-lo)/2
       if sumFloor(mid) <= T {
           lo = mid + 1
       } else {
           hi = mid
       }
   }
   x1 := lo
   // answer is x1 - x0
   fmt.Println(x1 - x0)
}
