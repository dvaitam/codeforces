package main

import (
   "fmt"
)

func main() {
   var n int64
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   // Maximum total kills over three turns is n + n/2 = 3n/2 for even n
   result := n*3/2
   fmt.Println(result)
}
