package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   infinite := false
   var ans int
   for i := 1; i < n; i++ {
       x, y := a[i-1], a[i]
       // triangle and square touch in a segment -> infinite
       if (x == 2 && y == 3) || (x == 3 && y == 2) {
           infinite = true
           break
       }
       // circle and triangle -> 3 points
       if (x == 1 && y == 2) || (x == 2 && y == 1) {
           ans += 3
       }
       // circle and square -> 4 points
       if (x == 1 && y == 3) || (x == 3 && y == 1) {
           ans += 4
       }
   }
   if infinite {
       fmt.Println("Infinite")
   } else {
       fmt.Println("Finite")
       fmt.Println(ans)
   }
}
