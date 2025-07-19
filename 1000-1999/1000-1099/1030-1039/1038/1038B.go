package main

import "fmt"

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   if n <= 2 {
       fmt.Println("No")
       return
   }
   fmt.Println("Yes")
   if n%2 == 0 {
       fmt.Printf("2 1 %d\n", n)
       fmt.Printf("%d", n-2)
       for i := 2; i < n; i++ {
           fmt.Printf(" %d", i)
       }
       fmt.Println()
   } else {
       k := (n + 1) / 2
       fmt.Printf("1 %d\n", k)
       fmt.Printf("%d", n-1)
       for i := 1; i <= n; i++ {
           if i == k {
               continue
           }
           fmt.Printf(" %d", i)
       }
       fmt.Println()
   }
}
