package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, t int
   var k int64
   if _, err := fmt.Fscan(in, &n, &t, &k); err != nil {
       return
   }
   a := make([]int64, n)
   b := make([]int64, n)
   c := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i], &b[i], &c[i])
   }
   // Estimate right bound for trains
   var sumd int64
   for i := 0; i < n; i++ {
       need := a[i] + int64(t)*b[i] - c[i]
       if need > 0 {
           sumd += need
       }
   }
   left := int64(0)
   right := (sumd + k - 1) / k
   // add t as buffer
   right += int64(t)

   // feasibility check for M trains
   feasible := func(M int64) bool {
       used := int64(0)
       B := make([]int64, n)
       copy(B, a)
       for h := 1; h <= t; h++ {
           // compute needed removals this hour
           prefix := make([]int64, n+1)
           for i := 0; i < n; i++ {
               prefix[i+1] = prefix[i] + B[i]
           }
           allZero := true
           var maxCap int64
           for i := 0; i < n; i++ {
               // must remove at least (B[i] + b[i] - c[i])
               req := B[i] + b[i] - c[i]
               if req < 0 {
                   req = 0
               } else {
                   allZero = false
               }
               needCap := prefix[i] + req
               if needCap > maxCap {
                   maxCap = needCap
               }
           }
           if !allZero {
               // number of trains this hour
               R := (maxCap + k - 1) / k
               if used+R > M {
                   return false
               }
               used += R
               rem := R * k
               // serve B in order
               for i := 0; i < n && rem > 0; i++ {
                   take := B[i]
                   if take > rem {
                       take = rem
                   }
                   B[i] -= take
                   rem -= take
               }
           }
           // arrivals at end of hour
           for i := 0; i < n; i++ {
               B[i] += b[i]
           }
       }
       return true
   }

   // binary search minimal M
   for left < right {
       mid := (left + right) / 2
       if feasible(mid) {
           right = mid
       } else {
           left = mid + 1
       }
   }
   // output answer
   fmt.Println(left)
}
