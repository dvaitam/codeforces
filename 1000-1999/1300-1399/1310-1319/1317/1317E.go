package main

import "fmt"

func main() {
   var n int
   fmt.Scan(&n)
   if n < 2 {
       fmt.Println("NO")
       return
   }
   for i := 2; i*i <= n; i++ {
       if n%i == 0 {
           fmt.Println("NO")
           return
       }
   }
   fmt.Println("YES")
}
