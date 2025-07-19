package main

import "fmt"

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&a[i])
   }
   if n == 0 {
       fmt.Println(-1)
       return
   }
   t := a[0]
   for i := 1; i < n; i++ {
       t = gcd(a[i], t)
   }
   if t != a[0] {
       fmt.Println(-1)
       return
   }
   total := 2*int64(n) - 1
   fmt.Println(total)
   fmt.Print(a[0])
   for i := 1; i < n; i++ {
       fmt.Printf(" %d %d", a[i], a[0])
   }
}
