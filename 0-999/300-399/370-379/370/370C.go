package main

import (
   "fmt"
   "sort"
)

func main() {
   var n, m int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&a[i])
   }
   sort.Ints(a)
   half := n / 2
   k := 0
   for i := 0; i < n; i++ {
       if a[i] != a[(i+half)%n] {
           k++
       }
   }
   fmt.Println(k)
   for i := 0; i < n; i++ {
       fmt.Println(a[i], a[(i+half)%n])
   }
}
