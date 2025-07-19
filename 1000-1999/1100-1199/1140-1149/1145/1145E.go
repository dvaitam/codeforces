package main

import (
   "fmt"
)

func main() {
   for i := 21; i < 51; i++ {
       // compute min(i, 25)
       m := i
       if 25 < i {
           m = 25
       }
       // evaluate condition
       a := (m + i) % (2 + i%3) > 0
       if a {
           fmt.Println(1)
       } else {
           fmt.Println(0)
       }
   }
}
