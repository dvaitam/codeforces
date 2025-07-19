package main

import (
   "fmt"
   "strings"
)

func main() {
   var n, k int
   if _, err := fmt.Scan(&n, &k); err != nil {
       return
   }
   fmt.Println("YES")
   blank := strings.Repeat(".", n)
   if k%2 == 0 {
       half := k / 2
       row := "." + strings.Repeat("#", half) + strings.Repeat(".", n-1-half)
       fmt.Println(blank)
       fmt.Println(row)
       fmt.Println(row)
       fmt.Println(blank)
   } else if k >= n-2 {
       fmt.Println(blank)
       row2 := "." + strings.Repeat("#", n-2) + "."
       rem := k - (n - 2)
       half := rem / 2
       midDots := (n - 2) - rem
       row3 := "." + strings.Repeat("#", half) + strings.Repeat(".", midDots) + strings.Repeat("#", half) + "."
       fmt.Println(row2)
       fmt.Println(row3)
       fmt.Println(blank)
   } else {
       fmt.Println(blank)
       left := (n - k) / 2
       row2 := strings.Repeat(".", left) + strings.Repeat("#", k) + strings.Repeat(".", left)
       fmt.Println(row2)
       fmt.Println(blank)
       fmt.Println(blank)
   }
}
