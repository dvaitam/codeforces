package main

import (
   "fmt"
)

func main() {
   var x int
   if _, err := fmt.Scan(&x); err != nil {
       return
   }
   s := fmt.Sprint(x)
   has := make([]bool, 10)
   for _, ch := range s {
       has[ch-'0'] = true
   }
   ans := 0
   for i := 1; i*i <= x; i++ {
       if x%i == 0 {
           d1 := i
           d2 := x / i
           if shares(d1, has) {
               ans++
           }
           if d2 != d1 && shares(d2, has) {
               ans++
           }
       }
   }
   fmt.Println(ans)
}

// shares reports whether d and x share at least one digit, where has marks digits in x
func shares(d int, has []bool) bool {
   if d == 0 {
       return has[0]
   }
   for d > 0 {
       if has[d%10] {
           return true
       }
       d /= 10
   }
   return false
}
