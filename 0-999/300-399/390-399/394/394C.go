package main

import "fmt"

func main() {
   // read dimensions

   var n, m int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   // counts: a[0]=11, a[1]=others (10/01), a[3]=00
   a := make([]int, 4)
   // count domino types
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           var s string
           fmt.Scan(&s)
           if s == "11" {
               a[0]++
           } else if s == "00" {
               a[3]++
           } else {
               a[1]++
           }
       }
   }
   ans := []string{"11", "10", "01", "00"}
   // build rows
   for i := 0; i < n; i++ {
       p := make([]int, m)
       for j := 0; j < m; j++ {
           if a[0] > 0 {
               p[j] = 0
               a[0]--
           } else if a[1] > 0 {
               p[j] = 1 + (i & 1)
               a[1]--
           } else {
               p[j] = 3
           }
       }
       // print row in snake order
       if i&1 == 1 {
           for j := 0; j < m; j++ {
               if j > 0 {
                   fmt.Print(" ")
               }
               fmt.Print(ans[p[j]])
           }
       } else {
           for k := 0; k < m; k++ {
               if k > 0 {
                   fmt.Print(" ")
               }
               j := m - 1 - k
               fmt.Print(ans[p[j]])
           }
       }
       fmt.Println()
   }
}
