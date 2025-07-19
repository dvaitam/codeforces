package main

import (
   "fmt"
   "strings"
)

func main() {
   var m string
   var a, b string
   if _, err := fmt.Scan(&m, &a, &b); err != nil {
       return
   }
   mByte := m[0]
   t := "6789TJQKA"
   win := false
   if (a[1] == mByte && b[1] != mByte) ||
      (a[1] == b[1] && strings.IndexByte(t, a[0]) > strings.IndexByte(t, b[0])) {
       win = true
   }
   if win {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
