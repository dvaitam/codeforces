package main

import (
   "fmt"
)

func main() {
   var n, m, k int64
   if _, err := fmt.Scan(&n, &m, &k); err != nil {
       return
   }
   // ensure n <= m for faster counting
   if n > m {
       n, m = m, n
   }
   var left, right, ans int64 = 1, n * m, 0
   for left <= right {
       mid := (left + right) / 2
       cnt := countLE(mid, n, m)
       if cnt >= k {
           ans = mid
           right = mid - 1
       } else {
           left = mid + 1
       }
   }
   fmt.Println(ans)
}

// countLE returns the number of elements in n x m multiplication table <= x
func countLE(x, n, m int64) (cnt int64) {
   for i := int64(1); i <= n; i++ {
       t := x / i
       if t > m {
           t = m
       }
       cnt += t
   }
   return
}
