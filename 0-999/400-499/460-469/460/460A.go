package main

import (
   "fmt"
)

func main() {
   var n, m int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   days := 0
   socks := n
   for socks > 0 {
       days++
       socks--
       if days % m == 0 {
           socks++
       }
   }
   fmt.Println(days)
}
