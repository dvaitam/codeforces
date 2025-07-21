package main

import (
   "fmt"
)

func main() {
   var m, s int
   if _, err := fmt.Scan(&m, &s); err != nil {
       return
   }
   // impossible cases
   if (s == 0 && m > 1) || s > 9*m {
       fmt.Println("-1 -1")
       return
   }
   // special case: zero
   if s == 0 && m == 1 {
       fmt.Println("0 0")
       return
   }
   // largest number
   rem := s
   max := make([]byte, m)
   for i := 0; i < m; i++ {
       d := 9
       if rem < 9 {
           d = rem
       }
       max[i] = byte('0' + d)
       rem -= d
   }
   // smallest number
   rem = s
   min := make([]byte, m)
   for i := 0; i < m; i++ {
       // lowest digit: can't be zero at first if more than one digit
       low := 0
       if i == 0 {
           low = 1
       }
       for d := low; d <= 9; d++ {
           // check if remainder achievable
           if rem-d < 0 {
               break
           }
           if rem-d <= 9*(m-i-1) {
               min[i] = byte('0' + d)
               rem -= d
               break
           }
       }
   }
   fmt.Printf("%s %s", string(min), string(max))
}
