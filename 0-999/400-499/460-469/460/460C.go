package main

import (
   "fmt"
)

func main() {
   var n, m, w int
   if _, err := fmt.Scan(&n, &m, &w); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&a[i])
   }
   // find minimum initial height
   minA := a[0]
   for _, v := range a {
       if v < minA {
           minA = v
       }
   }
   // binary search for maximum achievable height
   low := minA
   high := minA + int64(m) + 1
   var ans int64 = minA
   for low <= high {
       mid := (low + high) / 2
       if canReach(a, m, w, mid) {
           ans = mid
           low = mid + 1
       } else {
           high = mid - 1
       }
   }
   fmt.Println(ans)
}

// canReach checks if it's possible to make all flowers at least height h
// using at most m operations of watering w contiguous flowers by +1
func canReach(a []int64, m, w int, h int64) bool {
   n := len(a)
   ops := make([]int64, n)
   var added int64
   var used int64
   for i := 0; i < n; i++ {
       if i >= w {
           added -= ops[i-w]
       }
       curr := a[i] + added
       if curr < h {
           need := h - curr
           used += need
           if used > int64(m) {
               return false
           }
           added += need
           ops[i] = need
       }
   }
   return true
}
