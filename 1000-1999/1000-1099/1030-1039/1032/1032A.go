package main

import "fmt"

func main() {
   var n, k int
   if _, err := fmt.Scan(&n, &k); err != nil {
       return
   }
   cnt := make([]int, 101)
   for i := 0; i < n; i++ {
       var x int
       fmt.Scan(&x)
       if x >= 0 && x <= 100 {
           cnt[x]++
       }
   }
   mx := 0
   for i := 0; i <= 100; i++ {
       if cnt[i] > mx {
           mx = cnt[i]
       }
   }
   if k > 0 && mx%k != 0 {
       mx = (mx/k + 1) * k
   }
   ans := 0
   for i := 0; i <= 100; i++ {
       if cnt[i] > 0 {
           ans += mx - cnt[i]
       }
   }
   fmt.Print(ans)
}
