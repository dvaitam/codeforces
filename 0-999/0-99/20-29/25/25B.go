package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   res := make([]byte, 0, len(s)+len(s)/2)
   i := 0
   for i < n {
       rem := n - i
       if rem == 3 {
           res = append(res, s[i], s[i+1], s[i+2])
           break
       }
       if rem == 2 {
           res = append(res, s[i], s[i+1])
           break
       }
       // rem > 3: take two and a dash
       res = append(res, s[i], s[i+1], '-')
       i += 2
   }
   fmt.Println(string(res))
}
