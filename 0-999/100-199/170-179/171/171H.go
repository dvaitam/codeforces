package main

import (
   "fmt"
)

func main() {
   var a, b int
   if _, err := fmt.Scan(&a, &b); err != nil {
       return
   }
   q := b / a
   r := b % a
   fmt.Printf("%d %d", q, r)
}
