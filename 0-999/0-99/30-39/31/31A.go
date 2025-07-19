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
   for i := 0; i < n && !found; i++ {
       for j := i + 1; j < n && !found; j++ {
           for k := j + 1; k < n && !found; k++ {
               if a[i]+a[j] == a[k] {
                   fmt.Printf("%d %d %d\n", k+1, j+1, i+1)
                   found = true
               } else if a[i]+a[k] == a[j] {
                   fmt.Printf("%d %d %d\n", j+1, k+1, i+1)
                   found = true
               } else if a[j]+a[k] == a[i] {
                   fmt.Printf("%d %d %d\n", i+1, j+1, k+1)
                   found = true
               }
           }
       }
   }
   if !found {
       fmt.Println(-1)
   }
}
