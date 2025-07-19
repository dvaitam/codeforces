package main

import "fmt"

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&a[i])
   }
   // find first element different from the first
   idx := -1
   for i := 1; i < n; i++ {
       if a[i] != a[0] {
           idx = i
           break
       }
   }
   if idx == -1 {
       fmt.Println(0)
       return
   }
   // determine initial trend and look for breaking point
   if a[idx] > a[0] {
       for j := idx + 1; j < n; j++ {
           if a[j] < a[j-1] {
               fmt.Println(3)
               fmt.Printf("1 %d %d\n", j, j+1)
               return
           }
       }
   } else {
       for j := idx + 1; j < n; j++ {
           if a[j] > a[j-1] {
               fmt.Println(3)
               fmt.Printf("1 %d %d\n", j, j+1)
               return
           }
       }
   }
   fmt.Println(0)
}
