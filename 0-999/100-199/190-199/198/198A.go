package main

import (
   "fmt"
)

func main() {
   var k, b, n, t int64
   if _, err := fmt.Scan(&k, &b, &n, &t); err != nil {
       return
   }
   cur := int64(1)
   var steps int64
   for steps < n && k*cur+b <= t {
       cur = k*cur + b
       steps++
   }
   // Minimum seconds needed in second experiment
   result := n - steps
   fmt.Println(result)
}
