package main

import "fmt"

func main() {
   var n int
   var K int64
   fmt.Scan(&n, &K)
   var s1, s2 string
   fmt.Scan(&s1, &s2)

   var cnt, cur, ans int64 = 0, 1, 0
   // iterate over positions
   for i := 0; i < n; i++ {
       // double current possibilities
       cur <<= 1
       if s1[i] == 'b' {
           cur--
       }
       if s2[i] == 'a' {
           cur--
       }
       // if exceeds K, accumulate and break
       if cur > K {
           ans = cnt + K*int64(n-i)
           break
       }
       cnt += cur
       // optimistic add of remaining with same cur
       ans = cnt + cur*int64(n-1-i)
   }
   fmt.Println(ans)
}
