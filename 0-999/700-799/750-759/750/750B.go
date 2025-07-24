package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   lat := 0
   for i := 0; i < n; i++ {
       var t int
       var dir string
       fmt.Fscan(in, &t, &dir)
       // At poles, only allowed moves
       if (lat == 0 && dir != "South") || (lat == 20000 && dir != "North") {
           fmt.Println("NO")
           return
       }
       // Update latitude
       switch dir {
       case "South":
           lat += t
       case "North":
           lat -= t
       }
       // Check bounds
       if lat < 0 || lat > 20000 {
           fmt.Println("NO")
           return
       }
   }
   // Must end at North Pole
   if lat != 0 {
       fmt.Println("NO")
   } else {
       fmt.Println("YES")
   }
}
