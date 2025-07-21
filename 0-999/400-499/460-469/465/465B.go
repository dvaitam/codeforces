package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   unread := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&unread[i])
   }
   totalUnread, segments := 0, 0
   for i := 0; i < n; i++ {
       if unread[i] == 1 {
           totalUnread++
           if i == 0 || unread[i-1] == 0 {
               segments++
           }
       }
   }
   if totalUnread == 0 {
       fmt.Println(0)
   } else {
       fmt.Println(totalUnread + segments - 1)
   }
}
