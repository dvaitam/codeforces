package main

import (
   "fmt"
   "sort"
)

func main() {
   var n int
   var s string
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   a := make([]byte, n)
   b := make([]byte, n)
   for i := 0; i < n; i++ {
       a[i] = s[i]
       b[i] = s[n+i]
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
   less := true
   greater := true
   for i := 0; i < n; i++ {
       if a[i] >= b[i] {
           less = false
       }
       if a[i] <= b[i] {
           greater = false
       }
   }
   if less || greater {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
