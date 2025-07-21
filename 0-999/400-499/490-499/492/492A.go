package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   total := 0
   height := 0
   for i := 1; ; i++ {
       // cubes needed for level i: 1+2+...+i = i*(i+1)/2
       need := i * (i + 1) / 2
       if total+need > n {
           break
       }
       total += need
       height++
   }
   fmt.Println(height)
}
