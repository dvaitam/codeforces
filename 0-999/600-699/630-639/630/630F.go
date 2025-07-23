package main

import "fmt"

func comb(n, k int64) int64 {
   if k < 0 || k > n {
      return 0
   }
   if k > n-k {
      k = n - k
   }
   res := int64(1)
   for i := int64(1); i <= k; i++ {
      res = res * (n - i + 1) / i
   }
   return res
}

func main() {
   var n int64
   if _, err := fmt.Scan(&n); err != nil {
      return
   }
   ans := comb(n, 5) + comb(n, 6) + comb(n, 7)
   fmt.Println(ans)
}
