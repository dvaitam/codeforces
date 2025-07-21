package main

import (
   "fmt"
)

func main() {
   var a, b string
   if _, err := fmt.Scan(&a, &b); err != nil {
       return
   }
   cnt47, cnt74 := 0, 0
   for i := range a {
       if a[i] == '4' && b[i] == '7' {
           cnt47++
       } else if a[i] == '7' && b[i] == '4' {
           cnt74++
       }
   }
   if cnt47 > cnt74 {
       fmt.Println(cnt47)
   } else {
       fmt.Println(cnt74)
   }
}
