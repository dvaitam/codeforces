package main

import (
   "fmt"
)

func main() {
   var n, m int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   // Build the arrangement
   res := make([]byte, 0, n+m)
   // If more girls, start with a girl
   if m > n {
       res = append(res, 'G')
       m--
   }
   // Alternate as much as possible
   for n+m > 0 {
       if n > 0 {
           res = append(res, 'B')
           n--
       }
       if m > 0 {
           res = append(res, 'G')
           m--
       }
   }
   // Output result
   fmt.Println(string(res))
}
