package main

import (
   "fmt"
)

func main() {
   var x uint64
   if _, err := fmt.Scan(&x); err != nil {
       return
   }
   fmt.Println(x)
}
