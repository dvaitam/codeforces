package main

import "fmt"

func main() {
   // Read input values
   var n, m int64
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   current := int64(1)
   total := int64(0)
   for i := int64(0); i < m; i++ {
       var a int64
       if _, err := fmt.Scan(&a); err != nil {
           return
       }
       if a >= current {
           total += a - current
       } else {
           total += n - (current - a)
       }
       current = a
   }
   fmt.Println(total)
}
