package main

import (
   "fmt"
)

// check returns true if the subarray a[l..r] is sorted in non-decreasing order.
func check(a []int, l, r int) bool {
   for i := l + 1; i <= r; i++ {
       if a[i] < a[i-1] {
           return false
       }
   }
   return true
}

// thanos returns the size of the largest sorted subarray obtainable
// by recursively removing either the first or second half if unsorted.
func thanos(a []int, l, r int) int {
   if l == r {
       return 1
   }
   if check(a, l, r) {
       return r - l + 1
   }
   m := (l + r) / 2
   left := thanos(a, l, m)
   right := thanos(a, m+1, r)
   if left > right {
       return left
   }
   return right
}

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&a[i])
   }
   res := thanos(a, 0, n-1)
   fmt.Println(res)
}
