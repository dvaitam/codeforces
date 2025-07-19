package main

import "fmt"

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       if _, err := fmt.Scan(&a[i]); err != nil {
           return
       }
   }
   found := false
   // try all combinations of + and - for each angle
   for mask := 0; mask < (1 << n); mask++ {
       sum := 0
       for i := 0; i < n; i++ {
           if (mask>>i)&1 == 1 {
               sum += a[i]
           } else {
               sum -= a[i]
           }
       }
       if ((sum%360)+360)%360 == 0 {
           found = true
           break
       }
   }
   if found {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
