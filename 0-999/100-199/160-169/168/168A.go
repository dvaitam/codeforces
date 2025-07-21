package main

import (
   "fmt"
)

func main() {
   var n, x, y int64
   // Read total citizens n, number of wizards x, percentage y
   if _, err := fmt.Scan(&n, &x, &y); err != nil {
       return
   }
   // Compute minimum required demonstrators: ceil(n * y / 100)
   required := (n*y + 99) / 100
   // Clones needed beyond existing wizards
   clones := required - x
   if clones < 0 {
       clones = 0
   }
   fmt.Println(clones)
}
