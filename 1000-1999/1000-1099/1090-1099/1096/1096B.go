package main

import "fmt"

const MOD = 998244353

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }

   var l int64 = 1
   for i := 1; i < n; i++ {
       if s[i] == s[0] {
           l++
       } else {
           break
       }
   }
   var r int64 = 1
   for i := n - 2; i >= 0; i-- {
       if s[i] == s[n-1] {
           r++
       } else {
           break
       }
   }

   var res int64
   if s[0] == s[n-1] {
       res = (l * r) % MOD
   } else {
       res = (l + r - 1) % MOD
   }
   fmt.Print(res)
}
