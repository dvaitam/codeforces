package main

import (
   "fmt"
   "strings"
)

func main() {
   var a, b, x int
   if _, err := fmt.Scan(&a, &b, &x); err != nil {
       return
   }

   var sb0, sb1 strings.Builder
   for i := 0; i <= x; i++ {
       if i%2 == 0 {
           sb0.WriteByte('0')
           sb1.WriteByte('1')
       } else {
           sb0.WriteByte('1')
           sb1.WriteByte('0')
       }
   }

   alt0 := sb0.String()
   alt1 := sb1.String()

   if x%2 == 1 {
       atemp := a - (x+1)/2
       btemp := b - (x+1)/2
       if atemp >= 0 && btemp >= 0 {
           res := strings.Repeat("0", atemp) + alt0 + strings.Repeat("1", btemp)
           fmt.Println(res)
           return
       }
   } else {
       atemp := a - 1 - x/2
       btemp := b - x/2
       if atemp >= 0 && btemp >= 0 {
           prefix := alt0[:len(alt0)-1]
           suffix := alt0[len(alt0)-1:]
           res := strings.Repeat("0", atemp) + prefix + strings.Repeat("1", btemp) + suffix
           fmt.Println(res)
           return
       }

       atemp2 := a - x/2
       btemp2 := b - 1 - x/2
       if atemp2 >= 0 && btemp2 >= 0 {
           prefix := alt1[:len(alt1)-1]
           suffix := alt1[len(alt1)-1:]
           res := strings.Repeat("1", btemp2) + prefix + strings.Repeat("0", atemp2) + suffix
           fmt.Println(res)
           return
       }
   }
}
